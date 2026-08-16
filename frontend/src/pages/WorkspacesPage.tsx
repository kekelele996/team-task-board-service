import { useEffect, useState } from 'react';
import { Button, Card, Col, Empty, Form, Input, Modal, Row, Space, Typography, message } from 'antd';
import { LogoutOutlined, PlusOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { workspaceApi } from '../api';
import { errorMessage } from '../api/client';
import { useAuthStore } from '../store/auth';
import type { Workspace } from '../types';
import NotificationBell from '../components/NotificationBell';

export default function WorkspacesPage() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  const load = async () => {
    setLoading(true);
    try {
      setWorkspaces(await workspaceApi.list());
    } catch (error) {
      message.error(errorMessage(error));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const onCreate = async (values: { name: string; description?: string }) => {
    setSubmitting(true);
    try {
      const created = await workspaceApi.create({ name: values.name, description: values.description ?? '' });
      message.success('工作区已创建');
      setOpen(false);
      form.resetFields();
      navigate(`/workspaces/${created.id}/boards`);
    } catch (error) {
      message.error(errorMessage(error));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ maxWidth: 1080, margin: '0 auto', padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>
          我的工作区
        </Typography.Title>
        <Space>
          <NotificationBell />
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>
            新建工作区
          </Button>
          <Button
            icon={<LogoutOutlined />}
            onClick={() => {
              logout();
              navigate('/login');
            }}
          >
            退出 {user?.username}
          </Button>
        </Space>
      </div>
      {loading ? (
        <Card loading />
      ) : workspaces.length === 0 ? (
        <Empty description="还没有工作区，点击右上角新建" />
      ) : (
        <Row gutter={[16, 16]}>
          {workspaces.map((workspace) => (
            <Col span={8} key={workspace.id}>
              <Card
                hoverable
                onClick={() => navigate(`/workspaces/${workspace.id}/boards`)}
              >
                <Typography.Title level={5} style={{ marginTop: 0 }}>
                  {workspace.name}
                </Typography.Title>
                <Typography.Paragraph type="secondary" ellipsis={{ rows: 2 }}>
                  {workspace.description || '暂无描述'}
                </Typography.Paragraph>
                <Typography.Text type="secondary">所有者：{workspace.owner?.username ?? '—'}</Typography.Text>
              </Card>
            </Col>
          ))}
        </Row>
      )}
      <Modal
        title="新建工作区"
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical" onFinish={onCreate}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input placeholder="例如：产品研发" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="工作区用途" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
