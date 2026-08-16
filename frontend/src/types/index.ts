export interface User {
  id: number;
  username: string;
  email: string;
}

export type Role = 'admin' | 'editor' | 'viewer';
export type Priority = 'high' | 'medium' | 'low';

export interface TokenResponse {
  token: string;
  user: User;
}

export interface Workspace {
  id: number;
  name: string;
  description: string;
  owner_id: number;
  owner?: User;
  created_at: string;
  updated_at: string;
}

export interface Member {
  id: number;
  workspace_id: number;
  user_id: number;
  role: Role;
  username: string;
  email: string;
  created_at: string;
}

export interface Board {
  id: number;
  workspace_id: number;
  name: string;
  description: string;
  created_by: number;
  columns?: BoardColumn[];
  created_at: string;
  updated_at: string;
}

export interface BoardColumn {
  id: number;
  board_id: number;
  name: string;
  position: number;
  color?: string;
}

export interface Tag {
  id: number;
  workspace_id: number;
  name: string;
  color: string;
}

export interface Subtask {
  id: number;
  task_id: number;
  title: string;
  completed: boolean;
  position: number;
}

export interface Comment {
  id: number;
  task_id: number;
  user_id: number;
  content: string;
  created_at: string;
  user?: User;
}

export interface Attachment {
  id: number;
  task_id: number;
  file_name: string;
  file_path: string;
  size: number;
  content_type: string;
  uploaded_by: number;
  created_at: string;
}

export interface Task {
  id: number;
  board_id: number;
  column_id: number;
  title: string;
  description: string;
  assignee_id?: number | null;
  priority: Priority;
  due_date?: string | null;
  position: number;
  created_by: number;
  assignee?: User | null;
  column?: BoardColumn | null;
  subtasks?: Subtask[];
  tags?: Tag[];
  comments?: Comment[];
  attachments?: Attachment[];
  created_at: string;
  updated_at: string;
}

export interface ActivityLog {
  id: number;
  workspace_id: number;
  board_id: number;
  task_id?: number | null;
  user_id: number;
  action: string;
  detail: string;
  created_at: string;
  user?: User;
}

export interface Notification {
  id: number;
  user_id: number;
  actor_id?: number | null;
  type: string;
  content: string;
  read: boolean;
  created_at: string;
  actor?: User;
}

export interface Stats {
  columns: { column_id: number; name: string; count: number }[];
  assignees: { assignee_id: number; username: string; count: number }[];
  weekly_trend: { day: string; count: number }[];
  overdue_tasks: Task[];
}

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}
