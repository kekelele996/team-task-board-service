import { useEffect, useState } from 'react';
import { Button, Card, Form, Input, Modal, Popconfirm, Select, Table, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useParams } from 'react-router-dom';
import { workspaceApi } from '../api';
import { errorMessage } from '../api/client';
import { useAuthStore } from '../store/auth';
import type { Member, Role } from '../types';

const roleLabels: Record<Role, string> = {
  admin: '管理员',
  editor: '编辑者',
  viewer: '观察者',
};

export default function MembersPage() {
  const { workspaceId } = useParams();
  const workspaceNumber = Number(workspaceId);
  const currentUser = useAuthStore((state) => state.user);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  const me = members.find((member) => member.user_id === currentUser?.id);
  const isAdmin = me?.role === 'admin';

  const load = async () => {
    setLoading(true);
    try {
      setMembers(await workspaceApi.members(workspaceNumber));
    } catch (error) {
      message.error(errorMessage(error));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void load();
  }, [workspaceNumber]);

  const invite = async (values: { email?: string; username?: string; role: Role }) => {
    setSubmitting(true);
    try {
      await workspaceApi.addMember(workspaceNumber, {
        email: values.email,
        username: values.username,
        role: values.role,
      });
      message.success('成员已邀请');
      setOpen(false);
      form.resetFields();
      void load();
    } catch (error) {
      message.error(errorMessage(error));
    } finally {
      setSubmitting(false);
    }
  };

  const changeRole = async (member: Member, role: Role) => {
    try {
      await workspaceApi.updateMember(workspaceNumber, member.user_id, role);
      message.success('角色已更新');
      void load();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const removeMember = async (member: Member) => {
    try {
      await workspaceApi.removeMember(workspaceNumber, member.user_id);
      message.success('成员已移除');
      void load();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const columns = [
    { title: '用户', dataIndex: 'username', key: 'username' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: Role, record: Member) =>
        isAdmin && record.user_id !== currentUser?.id ? (
          <Select
            value={role}
            style={{ width: 120 }}
            onChange={(value: Role) => void changeRole(record, value)}
            options={[
              { value: 'admin', label: '管理员' },
              { value: 'editor', label: '编辑者' },
              { value: 'viewer', label: '观察者' },
            ]}
          />
        ) : (
          <Tag color={role === 'admin' ? 'gold' : role === 'editor' ? 'blue' : 'default'}>
            {roleLabels[role]}
          </Tag>
        ),
    },
    ...(isAdmin
      ? [
          {
            title: '操作',
            key: 'action',
            render: (_: unknown, record: Member) =>
              record.user_id !== currentUser?.id ? (
                <Popconfirm title="确定移除该成员？" onConfirm={() => void removeMember(record)}>
                  <Button danger size="small">移除</Button>
                </Popconfirm>
              ) : null,
          },
        ]
      : []),
  ];

  return (
    <Card
      title="工作区成员"
      extra={isAdmin ? <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>邀请成员</Button> : null}
    >
      <Table rowKey="user_id" columns={columns} dataSource={members} loading={loading} pagination={false} />
      <Modal
        title="邀请成员"
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical" onFinish={invite} initialValues={{ role: 'viewer' }}>
          <Form.Item name="email" label="邮箱（可选）" rules={[{ type: 'email' }]}>
            <Input placeholder="user@example.com" />
          </Form.Item>
          <Form.Item name="username" label="用户名（可选）">
            <Input placeholder="用户名" />
          </Form.Item>
          <Form.Item name="role" label="角色" rules={[{ required: true }]}>
            <Select
              options={[
                { value: 'admin', label: '管理员' },
                { value: 'editor', label: '编辑者' },
                { value: 'viewer', label: '观察者' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}
