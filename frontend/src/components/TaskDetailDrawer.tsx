import { useEffect, useMemo, useState } from 'react';
import type { UploadProps } from 'antd';
import {
  Button,
  Checkbox,
  Col,
  DatePicker,
  Drawer,
  Form,
  Input,
  List,
  Row,
  Select,
  Space,
  Tabs,
  Tag,
  Typography,
  Upload,
  message,
} from 'antd';
import { DeleteOutlined, InboxOutlined, PaperClipOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { taskApi, workspaceApi } from '../api';
import { errorMessage } from '../api/client';
import type { Attachment, Comment, Member, Role, Subtask, Tag as TagType, Task } from '../types';

interface Props {
  workspaceId: number;
  taskId: number;
  role: Role;
  onClose: () => void;
  onChanged: () => void;
}

const priorityOptions = [
  { value: 'high', label: '高' },
  { value: 'medium', label: '中' },
  { value: 'low', label: '低' },
];

export default function TaskDetailDrawer({ workspaceId, taskId, role, onClose, onChanged }: Props) {
  const [task, setTask] = useState<Task | null>(null);
  const [members, setMembers] = useState<Member[]>([]);
  const [tags, setTags] = useState<TagType[]>([]);
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);
  const [comment, setComment] = useState('');
  const [subtaskTitle, setSubtaskTitle] = useState('');
  const canEdit = role === 'admin' || role === 'editor';

  const load = async () => {
    try {
      const [detail, memberList, tagList] = await Promise.all([
        taskApi.get(workspaceId, taskId),
        workspaceApi.members(workspaceId),
        taskApi.tags(workspaceId),
      ]);
      setTask(detail);
      setMembers(memberList);
      setTags(tagList);
      form.setFieldsValue({
        title: detail.title,
        description: detail.description,
        assignee_id: detail.assignee_id,
        priority: detail.priority,
        due_date: detail.due_date ? dayjs(detail.due_date) : null,
        tag_ids: detail.tags?.map((tag) => tag.id) ?? [],
      });
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  useEffect(() => {
    void load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [taskId]);

  const save = async () => {
    const values = await form.validateFields();
    setSaving(true);
    try {
      const payload = {
        title: values.title,
        description: values.description ?? '',
        assignee_id: values.assignee_id ?? null,
        priority: values.priority,
        due_date: values.due_date ? values.due_date.format('YYYY-MM-DDTHH:mm:ssZ') : null,
        tag_ids: values.tag_ids ?? [],
      };
      await taskApi.update(workspaceId, taskId, payload);
      message.success('已保存');
      await load();
      onChanged();
    } catch (error) {
      message.error(errorMessage(error));
    } finally {
      setSaving(false);
    }
  };

  const addComment = async () => {
    if (!comment.trim()) return;
    try {
      await taskApi.createComment(workspaceId, taskId, comment.trim());
      setComment('');
      await load();
      onChanged();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const addSubtask = async () => {
    if (!subtaskTitle.trim()) return;
    try {
      await taskApi.createSubtask(workspaceId, taskId, subtaskTitle.trim());
      setSubtaskTitle('');
      await load();
      onChanged();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const toggleSubtask = async (subtask: Subtask) => {
    try {
      await taskApi.updateSubtask(workspaceId, subtask.id, { title: subtask.title, completed: !subtask.completed });
      await load();
      onChanged();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const removeSubtask = async (subtask: Subtask) => {
    try {
      await taskApi.removeSubtask(workspaceId, subtask.id);
      await load();
      onChanged();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const removeAttachment = async (attachment: Attachment) => {
    try {
      await taskApi.removeAttachment(workspaceId, attachment.id);
      await load();
      onChanged();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const uploadAttachment: NonNullable<UploadProps['customRequest']> = async (options) => {
    const file = options.file as unknown as File;
    try {
      await taskApi.uploadAttachment(workspaceId, taskId, file);
      options.onSuccess?.(undefined);
      await load();
      onChanged();
    } catch (error) {
      options.onError?.(error instanceof Error ? error : new Error(errorMessage(error)));
    }
  };

  const markdownContent = useMemo(() => task?.description || '', [task?.description]);

  const detailPane = (
    <Form form={form} layout="vertical" disabled={!canEdit}>
      <Form.Item name="title" label="标题" rules={[{ required: true }]}>
        <Input />
      </Form.Item>
      <Form.Item name="description" label="描述（支持 Markdown）">
        <Input.TextArea rows={5} />
      </Form.Item>
      <Row gutter={12}>
        <Col span={12}>
          <Form.Item name="assignee_id" label="负责人">
            <Select allowClear placeholder="选择成员" options={members.map((member) => ({ value: member.user_id, label: member.username }))} />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item name="priority" label="优先级">
            <Select options={priorityOptions} />
          </Form.Item>
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={12}>
          <Form.Item name="due_date" label="截止日期">
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item name="tag_ids" label="标签">
            <Select mode="multiple" allowClear placeholder="选择标签" options={tags.map((tag) => ({ value: tag.id, label: tag.name }))} />
          </Form.Item>
        </Col>
      </Row>
      {canEdit && (
        <Button type="primary" onClick={save} loading={saving}>
          保存修改
        </Button>
      )}
    </Form>
  );

  const descriptionPreview = task?.description ? (
    <div style={{ marginTop: 16 }}>
      <Typography.Text strong>描述预览</Typography.Text>
      <div style={{ background: '#fafafa', borderRadius: 8, padding: 12, marginTop: 8 }}>
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{markdownContent}</ReactMarkdown>
      </div>
    </div>
  ) : null;

  const subtaskPane = (
    <div>
      {canEdit && (
        <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
          <Input value={subtaskTitle} onChange={(e) => setSubtaskTitle(e.target.value)} placeholder="添加子任务" onPressEnter={addSubtask} />
          <Button type="primary" onClick={addSubtask}>添加</Button>
        </Space.Compact>
      )}
      <List
        dataSource={task?.subtasks ?? []}
        renderItem={(subtask) => (
          <List.Item
            actions={
              canEdit
                ? [
                    <Button key="delete" type="text" icon={<DeleteOutlined />} onClick={() => removeSubtask(subtask)} />,
                  ]
                : []
            }
          >
            <Checkbox checked={subtask.completed} disabled={!canEdit} onChange={() => toggleSubtask(subtask)}>
              <Typography.Text delete={subtask.completed}>{subtask.title}</Typography.Text>
            </Checkbox>
          </List.Item>
        )}
      />
    </div>
  );

  const commentPane = (
    <div>
      <List
        dataSource={task?.comments ?? []}
        locale={{ emptyText: '暂无评论' }}
        renderItem={(item: Comment) => (
          <List.Item>
            <List.Item.Meta
              title={item.user?.username ?? `用户 ${item.user_id}`}
              description={
                <>
                  <div>{item.content}</div>
                  <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                    {dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}
                  </Typography.Text>
                </>
              }
            />
          </List.Item>
        )}
      />
      <Space.Compact style={{ width: '100%', marginTop: 16 }}>
        <Input value={comment} onChange={(e) => setComment(e.target.value)} placeholder="写下评论..." onPressEnter={addComment} />
        <Button type="primary" onClick={addComment}>评论</Button>
      </Space.Compact>
    </div>
  );

  const attachmentPane = (
    <div>
      {canEdit && (
        <Upload.Dragger customRequest={uploadAttachment} showUploadList={false} multiple={false}>
          <p className="ant-upload-drag-icon"><InboxOutlined /></p>
          <p className="ant-upload-text">点击或拖拽文件上传</p>
        </Upload.Dragger>
      )}
      <List
        style={{ marginTop: 16 }}
        dataSource={task?.attachments ?? []}
        renderItem={(attachment) => (
          <List.Item
            actions={
              canEdit
                ? [<Button key="delete" type="text" icon={<DeleteOutlined />} onClick={() => removeAttachment(attachment)} />]
                : []
            }
          >
            <a href={attachment.file_path} target="_blank" rel="noreferrer">
              <PaperClipOutlined /> {attachment.file_name}
            </a>
            <Typography.Text type="secondary" style={{ fontSize: 12 }}>
              {attachment.size} bytes
            </Typography.Text>
          </List.Item>
        )}
      />
    </div>
  );

  return (
    <Drawer open width={560} onClose={onClose} title={task ? task.title : '任务详情'}>
      {task && (
        <>
          <Space wrap style={{ marginBottom: 12 }}>
            {task.tags?.map((tag) => (
              <Tag key={tag.id} color={tag.color}>{tag.name}</Tag>
            ))}
            {task.assignee && <Tag>负责人：{task.assignee.username}</Tag>}
          </Space>
          <Tabs
            defaultActiveKey="detail"
            items={[
              { key: 'detail', label: '详情', children: <>{detailPane}{descriptionPreview}</> },
              { key: 'subtasks', label: `子任务 (${task.subtasks?.length ?? 0})`, children: subtaskPane },
              { key: 'comments', label: `评论 (${task.comments?.length ?? 0})`, children: commentPane },
              { key: 'attachments', label: `附件 (${task.attachments?.length ?? 0})`, children: attachmentPane },
            ]}
          />
        </>
      )}
    </Drawer>
  );
}
