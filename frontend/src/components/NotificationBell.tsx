import { useEffect, useState } from 'react';
import { Badge, Button, Dropdown, Empty, Typography, message } from 'antd';
import { BellOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { notificationApi } from '../api';
import type { Notification } from '../types';

export default function NotificationBell() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [count, setCount] = useState(0);

  const load = async () => {
    try {
      const [items, unread] = await Promise.all([notificationApi.list(20), notificationApi.unreadCount()]);
      setNotifications(items);
      setCount(unread.count);
    } catch {
      // Keep the bell silent on transient failures.
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const markAll = async () => {
    try {
      await notificationApi.markAllRead();
      setCount(0);
      setNotifications((items) => items.map((item) => ({ ...item, read: true })));
      message.success('已全部标记为已读');
    } catch {
      // ignore
    }
  };

  const menu = {
    items: notifications.length
      ? [
          ...notifications.slice(0, 8).map((item) => ({
            key: String(item.id),
            label: (
              <div style={{ maxWidth: 320, whiteSpace: 'normal', padding: '4px 0' }}>
                <Typography.Text>{item.content}</Typography.Text>
                <div>
                  <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                    {dayjs(item.created_at).format('MM-DD HH:mm')}
                  </Typography.Text>
                </div>
              </div>
            ),
          })),
          { key: 'divider', type: 'divider' as const },
          { key: 'mark-all', label: '全部已读' },
        ]
      : [{ key: 'empty', label: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无通知" /> }],
    onClick: ({ key }: { key: string }) => {
      if (key === 'mark-all') void markAll();
    },
  };

  return (
    <Dropdown menu={menu} trigger={['click']} placement="bottomRight">
      <Button icon={<Badge count={count} size="small"><BellOutlined /></Badge>} type="text" />
    </Dropdown>
  );
}
