import { create } from 'zustand';
import type { Board, Member, Workspace } from '../types';

interface WorkspaceState {
  currentWorkspace: Workspace | null;
  currentBoard: Board | null;
  members: Member[];
  setCurrentWorkspace: (workspace: Workspace | null) => void;
  setCurrentBoard: (board: Board | null) => void;
  setMembers: (members: Member[]) => void;
  reset: () => void;
}

export const useWorkspaceStore = create<WorkspaceState>((set) => ({
  currentWorkspace: null,
  currentBoard: null,
  members: [],
  setCurrentWorkspace: (workspace) => set({ currentWorkspace: workspace }),
  setCurrentBoard: (board) => set({ currentBoard: board }),
  setMembers: (members) => set({ members }),
  reset: () => set({ currentWorkspace: null, currentBoard: null, members: [] }),
}));
