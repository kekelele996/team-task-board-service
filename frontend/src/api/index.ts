import { apiDelete, apiGet, apiPatch, apiPost, apiPut } from './client';
import type {
  ActivityLog,
  Attachment,
  Board,
  BoardColumn,
  Comment,
  Member,
  Notification,
  Stats,
  Subtask,
  Tag,
  Task,
  TokenResponse,
  User,
  Workspace,
} from '../types';

export const authApi = {
  register: (payload: { username: string; email: string; password: string }) =>
    apiPost<TokenResponse>('/auth/register', payload),
  login: (payload: { email: string; password: string }) =>
    apiPost<TokenResponse>('/auth/login', payload),
  me: () => apiGet<User>('/me'),
};

export const workspaceApi = {
  list: () => apiGet<Workspace[]>('/workspaces'),
  create: (payload: { name: string; description: string }) =>
    apiPost<Workspace>('/workspaces', payload),
  get: (id: number) => apiGet<Workspace>(`/workspaces/${id}`),
  update: (id: number, payload: { name: string; description: string }) =>
    apiPut<Workspace>(`/workspaces/${id}`, payload),
  remove: (id: number) => apiDelete<{ deleted: boolean }>(`/workspaces/${id}`),
  members: (id: number) => apiGet<Member[]>(`/workspaces/${id}/members`),
  addMember: (id: number, payload: { email?: string; username?: string; role: string }) =>
    apiPost<{ invited: boolean }>(`/workspaces/${id}/members`, payload),
  updateMember: (id: number, userId: number, role: string) =>
    apiPut<{ updated: boolean }>(`/workspaces/${id}/members/${userId}`, { role }),
  removeMember: (id: number, userId: number) =>
    apiDelete<{ removed: boolean }>(`/workspaces/${id}/members/${userId}`),
};

export const boardApi = {
  list: (workspaceId: number) => apiGet<Board[]>(`/workspaces/${workspaceId}/boards`),
  create: (workspaceId: number, payload: { name: string; description: string }) =>
    apiPost<Board>(`/workspaces/${workspaceId}/boards`, payload),
  get: (workspaceId: number, boardId: number) =>
    apiGet<Board>(`/workspaces/${workspaceId}/boards/${boardId}`),
  update: (workspaceId: number, boardId: number, payload: { name: string; description: string }) =>
    apiPut<Board>(`/workspaces/${workspaceId}/boards/${boardId}`, payload),
  remove: (workspaceId: number, boardId: number) =>
    apiDelete<{ deleted: boolean }>(`/workspaces/${workspaceId}/boards/${boardId}`),
  columns: (workspaceId: number, boardId: number) =>
    apiGet<BoardColumn[]>(`/workspaces/${workspaceId}/boards/${boardId}/columns`),
  createColumn: (workspaceId: number, boardId: number, payload: { name: string; color: string }) =>
    apiPost<BoardColumn>(`/workspaces/${workspaceId}/boards/${boardId}/columns`, payload),
  updateColumn: (workspaceId: number, columnId: number, payload: { name: string; color: string }) =>
    apiPut<BoardColumn>(`/workspaces/${workspaceId}/columns/${columnId}`, payload),
  removeColumn: (workspaceId: number, columnId: number) =>
    apiDelete<{ deleted: boolean }>(`/workspaces/${workspaceId}/columns/${columnId}`),
};

export const taskApi = {
  list: (workspaceId: number, boardId: number, params: Record<string, unknown>) =>
    apiGet<Task[]>(`/workspaces/${workspaceId}/boards/${boardId}/tasks`, params),
  create: (workspaceId: number, boardId: number, payload: Record<string, unknown>) =>
    apiPost<Task>(`/workspaces/${workspaceId}/boards/${boardId}/tasks`, payload),
  get: (workspaceId: number, taskId: number) =>
    apiGet<Task>(`/workspaces/${workspaceId}/tasks/${taskId}`),
  update: (workspaceId: number, taskId: number, payload: Record<string, unknown>) =>
    apiPut<Task>(`/workspaces/${workspaceId}/tasks/${taskId}`, payload),
  remove: (workspaceId: number, taskId: number) =>
    apiDelete<{ deleted: boolean }>(`/workspaces/${workspaceId}/tasks/${taskId}`),
  move: (workspaceId: number, taskId: number, payload: { column_id: number; position: number; task_ids?: number[] }) =>
    apiPatch<Task>(`/workspaces/${workspaceId}/tasks/${taskId}/move`, payload),
  createSubtask: (workspaceId: number, taskId: number, title: string) =>
    apiPost<Subtask>(`/workspaces/${workspaceId}/tasks/${taskId}/subtasks`, { title }),
  updateSubtask: (workspaceId: number, subtaskId: number, payload: { title: string; completed: boolean }) =>
    apiPut<Subtask>(`/workspaces/${workspaceId}/subtasks/${subtaskId}`, payload),
  removeSubtask: (workspaceId: number, subtaskId: number) =>
    apiDelete<{ deleted: boolean }>(`/workspaces/${workspaceId}/subtasks/${subtaskId}`),
  createComment: (workspaceId: number, taskId: number, content: string) =>
    apiPost<Comment>(`/workspaces/${workspaceId}/tasks/${taskId}/comments`, { content }),
  removeAttachment: (workspaceId: number, attachmentId: number) =>
    apiDelete<{ deleted: boolean }>(`/workspaces/${workspaceId}/attachments/${attachmentId}`),
  tags: (workspaceId: number) => apiGet<Tag[]>(`/workspaces/${workspaceId}/tags`),
  createTag: (workspaceId: number, payload: { name: string; color: string }) =>
    apiPost<Tag>(`/workspaces/${workspaceId}/tags`, payload),
  removeTag: (workspaceId: number, tagId: number) =>
    apiDelete<{ deleted: boolean }>(`/workspaces/${workspaceId}/tags/${tagId}`),
  activities: (workspaceId: number, boardId: number) =>
    apiGet<ActivityLog[]>(`/workspaces/${workspaceId}/boards/${boardId}/activities`),
  workspaceActivities: (workspaceId: number) =>
    apiGet<ActivityLog[]>(`/workspaces/${workspaceId}/activities`),
  stats: (workspaceId: number) => apiGet<Stats>(`/workspaces/${workspaceId}/stats`),
  uploadAttachment: (workspaceId: number, taskId: number, file: File) => {
    const form = new FormData();
    form.append('file', file);
    return apiPost<Attachment>(`/workspaces/${workspaceId}/tasks/${taskId}/attachments`, form);
  },
};

export const notificationApi = {
  list: (limit = 50) => apiGet<Notification[]>('/notifications', { limit }),
  unreadCount: () => apiGet<{ count: number }>('/notifications/unread-count'),
  markRead: (id: number) => apiPut<{ read: boolean }>(`/notifications/${id}/read`),
  markAllRead: () => apiPut<{ read: boolean }>('/notifications/read-all'),
};
