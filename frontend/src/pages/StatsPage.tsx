import { useEffect, useState } from 'react';
import { Card, Col, Empty, List, Row, Space, Spin, Typography, message } from 'antd';
import { useParams } from 'react-router-dom';
import dayjs from 'dayjs';
import { taskApi } from '../api';
import { errorMessage } from '../api/client';
import type { Stats } from '../types';

const palette = ['#1677ff', '#52c41a', '#faad14', '#eb2f96', '#13c2c2', '#722ed1', '#f5222d'];

function Donut({ data }: { data: { label: string; value: number }[] }) {
  const total = data.reduce((sum, item) => sum + item.value, 0);
  let offset = 0;
  const radius = 60;
  const circumference = 2 * Math.PI * radius;
  return (
    <svg width="180" height="140" viewBox="0 0 160 140">
      <circle cx="70" cy="70" r={radius} fill="none" stroke="#f0f0f0" strokeWidth="22" />
      {data.map((item, index) => {
        const dash = total ? (item.value / total) * circumference : 0;
        const node = (
          <circle
            key={item.label}
            cx="70"
            cy="70"
            r={radius}
            fill="none"
            stroke={palette[index % palette.length]}
            strokeWidth="22"
            strokeDasharray={`${dash} ${circumference - dash}`}
            strokeDashoffset={-offset}
            transform="rotate(-90 70 70)"
          />
        );
        offset += dash;
        return node;
      })}
    </svg>
  );
}

function Trend({ data }: { data: { day: string; count: number }[] }) {
  const max = Math.max(1, ...data.map((item) => item.count));
  const width = 480;
  const height = 180;
  const padding = 30;
  const points = data.map((item, index) => {
    const x = padding + (index * (width - padding * 2)) / Math.max(1, data.length - 1);
    const y = height - padding - (item.count / max) * (height - padding * 2);
    return { x, y, label: dayjs(item.day).format('MM-DD'), count: item.count };
  });
  return (
    <svg width="100%" height="200" viewBox={`0 0 ${width} ${height}`}>
      {points.map((point) => (
        <g key={point.label}>
          <line x1="0" y1={point.y} x2={width} y2={point.y} stroke="#f0f0f0" strokeDasharray="4" />
          <circle cx={point.x} cy={point.y} r="4" fill="#1677ff" />
          <text x={point.x} y={height - 10} textAnchor="middle" fontSize="11">{point.label}</text>
          <text x={point.x} y={point.y - 10} textAnchor="middle" fontSize="11">{point.count}</text>
        </g>
      ))}
      {points.length > 1 && (
        <polyline
          points={points.map((point) => `${point.x},${point.y}`).join(' ')}
          fill="none"
          stroke="#1677ff"
          strokeWidth="2"
        />
      )}
    </svg>
  );
}

export default function StatsPage() {
  const { workspaceId } = useParams();
  const workspaceNumber = Number(workspaceId);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        setStats(await taskApi.stats(workspaceNumber));
      } catch (error) {
        message.error(errorMessage(error));
      } finally {
        setLoading(false);
      }
    }
    void load();
  }, [workspaceNumber]);

  if (loading) {
    return <div style={{ display: 'flex', justifyContent: 'center', padding: 80 }}><Spin size="large" /></div>;
  }
  if (!stats) return <Empty />;

  const memberData = (stats.assignees ?? []).map((item) => ({ label: item.username, value: item.count }));
  const trendData = stats.weekly_trend ?? [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Row gutter={16}>
        {(stats.columns ?? []).map((column) => (
          <Col span={6} key={column.column_id}>
            <Card>
              <Typography.Text type="secondary">{column.name}</Typography.Text>
              <div style={{ fontSize: 32, fontWeight: 700 }}>{column.count}</div>
            </Card>
          </Col>
        ))}
      </Row>
      <Row gutter={16}>
        <Col span={10}>
          <Card title="成员任务分布">
            <Donut data={memberData} />
            <Space direction="vertical" style={{ width: '100%', marginTop: 8 }}>
              {memberData.map((item, index) => (
                <div key={item.label} style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <span style={{ width: 10, height: 10, borderRadius: 10, background: palette[index % palette.length], display: 'inline-block' }} />
                  <span>{item.label}</span>
                  <span style={{ marginLeft: 'auto' }}>{item.value}</span>
                </div>
              ))}
            </Space>
          </Card>
        </Col>
        <Col span={14}>
          <Card title="本周完成任务趋势">
            {trendData.length ? <Trend data={trendData} /> : <Empty description="本周暂无完成记录" />}
          </Card>
        </Col>
      </Row>
      <Card title="逾期任务">
        <List
          dataSource={stats.overdue_tasks ?? []}
          locale={{ emptyText: '没有逾期任务' }}
          renderItem={(task) => (
            <List.Item>
              <Space>
                <span>{task.title}</span>
                {task.due_date && <Typography.Text type="danger">{dayjs(task.due_date).format('YYYY-MM-DD')}</Typography.Text>}
                {task.assignee && <Typography.Text type="secondary">{task.assignee.username}</Typography.Text>}
              </Space>
            </List.Item>
          )}
        />
      </Card>
    </Space>
  );
}
