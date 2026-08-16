import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Button,
  Col,
  DatePicker,
  Empty,
  Form,
  Input,
  Modal,
  Row,
  Select,
  Space,
  Spin,
  Tag,
  message,
} from 'antd';
import { FilterOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { DragDropContext, Draggable, Droppable, type DropResult } from '@hello-pangea/dnd';
import dayjs from 'dayjs';
import { useOutletContext, useParams } from 'react-router-dom';
import { boardApi, taskApi, workspaceApi } from '../api';
import { errorMessage } from '../api/client';
import type { Board, BoardColumn, Member, Role, Tag as TagType, Task } from '../types';
import TaskCard from '../components/TaskCard';
import TaskDetailDrawer from '../components/TaskDetailDrawer';

interface OutletContext {
  workspaceId: number;
  role: Role;
  boards: Board[];
  setCurrentBoard: (board: Board) => void;
  refreshBoards: () => Promise<Board[]>;
}

const tagColors = ['magenta', 'red', 'volcano', 'orange', 'gold', 'lime', 'green', 'cyan', 'blue', 'geekblue', 'purple'];

export default function BoardPage() {
  const outlet = useOutletContext<OutletContext>();
  const { workspaceId, role, boards, setCurrentBoard } = outlet;
  const canEdit = role === 'admin' || role === 'editor';
  const boardId = Number(useParams().boardId || boards[0]?.id || 0);

  const [board, setBoard] = useState<Board | null>(null);
  const [columns, setColumns] = useState<BoardColumn[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [tags, setTags] = useState<TagType[]>([]);
  const [loading, setLoading] = useState(true);
  const [taskModalOpen, setTaskModalOpen] = useState(false);
  const [columnModalOpen, setColumnModalOpen] = useState(false);
  const [selectedTaskId, setSelectedTaskId] = useState<number | null>(null);
  const [keyword, setKeyword] = useState('');
  const [filterAssignee, setFilterAssignee] = useState<number | null>(null);
  const [filterPriority, setFilterPriority] = useState<string | null>(null);
  const [filterTag, setFilterTag] = useState<number | null>(null);
  const [filterDue, setFilterDue] = useState<string | null>(null);
  const [taskForm] = Form.useForm();
  const [columnForm] = Form.useForm();

  const loadBoard = useCallback(async () => {
    if (!boardId) return;
    setLoading(true);
    try {
      const [boardDetail, columnList, taskList, memberList, tagList] = await Promise.all([
        boardApi.get(workspaceId, boardId),
        boardApi.columns(workspaceId, boardId),
        taskApi.list(workspaceId, boardId, {}),
        workspaceApi.members(workspaceId),
        taskApi.tags(workspaceId),
      ]);
      setBoard(boardDetail);
      setColumns(columnList);
      setTasks(taskList);
      setMembers(memberList);
      setTags(tagList);
      setCurrentBoard(boardDetail);
    } catch (error) {
      message.error(errorMessage(error));
    } finally {
      setLoading(false);
    }
  }, [boardId, workspaceId, setCurrentBoard]);

  useEffect(() => {
    void loadBoard();
  }, [loadBoard]);

  const matches = useMemo(() => {
    const active = Boolean(keyword || filterAssignee || filterPriority || filterTag || filterDue);
    if (!active) return new Set<number>();
    return new Set(
      tasks
        .filter((task) => {
          if (filterAssignee && task.assignee_id !== filterAssignee) return false;
          if (filterPriority && task.priority !== filterPriority) return false;
          if (filterTag && !(task.tags ?? []).some((tag) => tag.id === filterTag)) return false;
          if (filterDue && task.due_date && dayjs(task.due_date).isAfter(filterDue, 'day')) return false;
          if (keyword) {
            const text = `${task.title} ${task.description ?? ''}`.toLowerCase();
            if (!text.includes(keyword.toLowerCase())) return false;
          }
          return true;
        })
        .map((task) => task.id),
    );
  }, [tasks, keyword, filterAssignee, filterPriority, filterTag, filterDue]);

  const hasActiveFilter = keyword !== '' || filterAssignee !== null || filterPriority !== null || filterTag !== null || filterDue !== null;

  const grouped = useMemo(() => {
    const map = new Map<number, Task[]>();
    for (const column of columns) map.set(column.id, []);
    for (const task of tasks) {
      const list = map.get(task.column_id) ?? [];
      list.push(task);
      map.set(task.column_id, list);
    }
    for (const list of map.values()) {
      list.sort((a, b) => a.position - b.position || a.id - b.id);
    }
    return map;
  }, [columns, tasks]);

  const onDragEnd = async (result: DropResult) => {
    if (!result.destination || !canEdit) return;
    const sourceColumnId = Number(result.source.droppableId);
    const destinationColumnId = Number(result.destination.droppableId);
    const taskId = Number(result.draggableId);

    const sourceTasks = Array.from(grouped.get(sourceColumnId) ?? []);
    const destinationTasks = Array.from(grouped.get(destinationColumnId) ?? []);
    const [moved] = sourceTasks.splice(result.source.index, 1);
    if (!moved) return;
    if (sourceColumnId === destinationColumnId) {
      destinationTasks.splice(result.destination.index, 0, moved);
    } else {
      destinationTasks.splice(result.destination.index, 0, moved);
    }

    const updated = tasks.map((task) => {
      if (task.id === taskId) return { ...task, column_id: destinationColumnId };
      return task;
    });
    setTasks(updated);

    const taskIds = destinationTasks.map((task) => task.id);
    try {
      await taskApi.move(workspaceId, taskId, {
        column_id: destinationColumnId,
        position: result.destination.index,
        task_ids: taskIds,
      });
      message.success('任务已移动');
    } catch (error) {
      message.error(errorMessage(error));
      void loadBoard();
    }
  };

  const createTask = async (values: { title: string; description?: string; column_id?: number; priority?: string }) => {
    try {
      await taskApi.create(workspaceId, boardId, {
        title: values.title,
        description: values.description ?? '',
        column_id: values.column_id ?? columns[0]?.id,
        priority: values.priority ?? 'medium',
      });
      message.success('任务已创建');
      setTaskModalOpen(false);
      taskForm.resetFields();
      void loadBoard();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  const createColumn = async (values: { name: string; color?: string }) => {
    try {
      await boardApi.createColumn(workspaceId, boardId, {
        name: values.name,
        color: values.color ?? tagColors[Math.floor(Math.random() * tagColors.length)],
      });
      message.success('列已创建');
      setColumnModalOpen(false);
      columnForm.resetFields();
      void loadBoard();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  if (loading && !board) {
    return <div style={{ display: 'flex', justifyContent: 'center', padding: 80 }}><Spin size="large" /></div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', gap: 12, marginBottom: 16, flexWrap: 'wrap', alignItems: 'center' }}>
        <Input
          allowClear
          prefix={<SearchOutlined />}
          placeholder="搜索任务标题或描述"
          style={{ width: 240 }}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
        />
        <Select
          allowClear
          placeholder="按负责人"
          style={{ width: 140 }}
          value={filterAssignee}
          onChange={setFilterAssignee}
          options={members.map((member) => ({ value: member.user_id, label: member.username }))}
        />
        <Select
          allowClear
          placeholder="按优先级"
          style={{ width: 130 }}
          value={filterPriority}
          onChange={setFilterPriority}
          options={[
            { value: 'high', label: '高' },
            { value: 'medium', label: '中' },
            { value: 'low', label: '低' },
          ]}
        />
        <Select
          allowClear
          placeholder="按标签"
          style={{ width: 140 }}
          value={filterTag}
          onChange={setFilterTag}
          options={tags.map((tag) => ({ value: tag.id, label: tag.name }))}
        />
        <DatePicker
          value={filterDue ? dayjs(filterDue) : null}
          onChange={(date) => setFilterDue(date ? date.format('YYYY-MM-DD') : null)}
          placeholder="截止日期"
        />
        <Space>
          {canEdit && (
            <>
              <Button icon={<PlusOutlined />} onClick={() => setTaskModalOpen(true)}>新建任务</Button>
              <Button icon={<FilterOutlined />} onClick={() => setColumnModalOpen(true)}>新增列</Button>
            </>
          )}
        </Space>
      </div>

      {columns.length === 0 ? (
        <Empty description="还没有列，点击“新增列”开始" />
      ) : (
        <DragDropContext onDragEnd={(result) => void onDragEnd(result)}>
          <div style={{ display: 'flex', gap: 16, overflowX: 'auto', paddingBottom: 12 }}>
            {columns.map((column) => (
              <Droppable droppableId={String(column.id)} key={column.id} isDropDisabled={!canEdit}>
                {(provided) => (
                  <div className="kanban-column" ref={provided.innerRef} {...provided.droppableProps}>
                    <div className="kanban-column-header">
                      <span>{column.name}</span>
                      <Tag>{grouped.get(column.id)?.length ?? 0}</Tag>
                    </div>
                    <div style={{ flex: 1, overflowY: 'auto', padding: '0 10px 10px' }}>
                      {(grouped.get(column.id) ?? []).map((task, index) => (
                        <Draggable draggableId={String(task.id)} index={index} key={task.id} isDragDisabled={!canEdit}>
                          {(dragProvided) => (
                            <div
                              ref={dragProvided.innerRef}
                              {...dragProvided.draggableProps}
                              {...dragProvided.dragHandleProps}
                              onClick={() => setSelectedTaskId(task.id)}
                            >
                              <TaskCard
                                task={{ ...task, column: column }}
                                dimmed={hasActiveFilter && !matches.has(task.id)}
                              />
                            </div>
                          )}
                        </Draggable>
                      ))}
                      {provided.placeholder}
                    </div>
                  </div>
                )}
              </Droppable>
            ))}
          </div>
        </DragDropContext>
      )}

      <Modal
        title="新建任务"
        open={taskModalOpen}
        onCancel={() => setTaskModalOpen(false)}
        onOk={() => taskForm.submit()}
      >
        <Form form={taskForm} layout="vertical" onFinish={createTask}>
          <Form.Item name="title" label="标题" rules={[{ required: true }]}>
            <Input placeholder="任务标题" />
          </Form.Item>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="column_id" label="所在列">
                <Select placeholder="选择列" options={columns.map((column) => ({ value: column.id, label: column.name }))} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="priority" label="优先级" initialValue="medium">
                <Select options={[
                  { value: 'high', label: '高' },
                  { value: 'medium', label: '中' },
                  { value: 'low', label: '低' },
                ]} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="新增列"
        open={columnModalOpen}
        onCancel={() => setColumnModalOpen(false)}
        onOk={() => columnForm.submit()}
      >
        <Form form={columnForm} layout="vertical" onFinish={createColumn}>
          <Form.Item name="name" label="列名" rules={[{ required: true }]}>
            <Input placeholder="例如：审核中" />
          </Form.Item>
        </Form>
      </Modal>

      {selectedTaskId && (
        <TaskDetailDrawer
          workspaceId={workspaceId}
          taskId={selectedTaskId}
          role={role}
          onClose={() => setSelectedTaskId(null)}
          onChanged={() => void loadBoard()}
        />
      )}
    </div>
  );
}
