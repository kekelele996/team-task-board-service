import { Avatar, Space, Tag, Tooltip, Typography } from 'antd';
import { CalendarOutlined, CommentOutlined, PaperClipOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import type { Task } from '../types';

const priorityMap = {
  high: { label: '高', className: 'priority-high' },
  medium: { label: '中', className: 'priority-medium' },
  low: { label: '低', className: 'priority-low' },
};

export default function TaskCard({ task, dimmed }: { task: Task; dimmed: boolean }) {
  const priority = priorityMap[task.priority] ?? priorityMap.medium;
  const overdue = task.due_date && dayjs(task.due_date).isBefore(dayjs(), 'day') && task.column?.name !== 'Done';
  return (
    <div className={`kanban-card ${dimmed ? 'dimmed' : ''}`}>
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8, marginBottom: 6 }}>
        <Typography.Text strong ellipsis={{ tooltip: task.title }}>
          {task.title}
        </Typography.Text>
        <span className={priority.className} style={{ borderRadius: 4, padding: '0 6px', fontSize: 12 }}>
          {priority.label}
        </span>
      </div>
      {task.tags && task.tags.length > 0 && (
        <div style={{ marginBottom: 6 }}>
          {task.tags.map((tag) => (
            <Tag key={tag.id} color={tag.color} style={{ marginBottom: 4 }}>
              {tag.name}
            </Tag>
          ))}
        </div>
      )}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', color: '#8c8c8c' }}>
        <Space size={8}>
          {task.assignee ? (
            <Tooltip title={task.assignee.username}>
              <Avatar size={20}>{task.assignee.username.slice(0, 1).toUpperCase()}</Avatar>
            </Tooltip>
          ) : null}
          {task.due_date && (
            <Tooltip title={`截止 ${dayjs(task.due_date).format('YYYY-MM-DD')}`}>
              <span style={{ color: overdue ? '#cf1322' : undefined }}>
                <CalendarOutlined /> {dayjs(task.due_date).format('MM-DD')}
              </span>
            </Tooltip>
          )}
        </Space>
        <Space size={10}>
          {task.comments && task.comments.length > 0 && (
            <span><CommentOutlined /> {task.comments.length}</span>
          )}
          {task.attachments && task.attachments.length > 0 && (
            <span><PaperClipOutlined /> {task.attachments.length}</span>
          )}
        </Space>
      </div>
    </div>
  );
}
