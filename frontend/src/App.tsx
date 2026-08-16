import type { ReactElement } from 'react';
import { Navigate, Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import { useAuthStore } from './store/auth';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import WorkspacesPage from './pages/WorkspacesPage';
import WorkspaceLayout from './pages/WorkspaceLayout';
import BoardPage from './pages/BoardPage';
import StatsPage from './pages/StatsPage';
import MembersPage from './pages/MembersPage';

function Protected({ children }: { children: ReactElement }) {
  const token = useAuthStore((state) => state.token);
  if (!token) return <Navigate to="/login" replace />;
  return children;
}

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          path="/"
          element={
            <Protected>
              <WorkspacesPage />
            </Protected>
          }
        />
        <Route
          path="/workspaces/:workspaceId"
          element={
            <Protected>
              <WorkspaceLayout />
            </Protected>
          }
        >
          <Route index element={<Navigate to="boards" replace />} />
          <Route path="boards" element={<BoardPage />} />
          <Route path="boards/:boardId" element={<BoardPage />} />
          <Route path="stats" element={<StatsPage />} />
          <Route path="members" element={<MembersPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}
