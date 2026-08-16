import { useEffect, useState } from 'react';
import { Layout, Menu, Spin, Typography, message } from 'antd';
import { AppstoreOutlined, BarChartOutlined, TeamOutlined } from '@ant-design/icons';
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom';
import { boardApi, workspaceApi } from '../api';
import { errorMessage } from '../api/client';
import { useAuthStore } from '../store/auth';
import { useWorkspaceStore } from '../store/workspace';
import type { Board, Member, Role } from '../types';
import NotificationBell from '../components/NotificationBell';

const { Sider, Content, Header } = Layout;

export default function WorkspaceLayout() {
  const { workspaceId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const user = useAuthStore((state) => state.user);
  const { currentWorkspace, currentBoard, setCurrentWorkspace, setCurrentBoard, setMembers } = useWorkspaceStore();
  const [boards, setBoards] = useState<Board[]>([]);
  const [role, setRole] = useState<Role>('viewer');
  const [loading, setLoading] = useState(true);

  const workspaceNumber = Number(workspaceId);

  const loadBoards = async () => {
    const boardList = await boardApi.list(workspaceNumber);
    setBoards(boardList);
    return boardList;
  };

  useEffect(() => {
    let cancelled = false;
    async function load() {
      if (!workspaceNumber) return;
      setLoading(true);
      try {
        const [workspace, boardList, memberList] = await Promise.all([
          workspaceApi.get(workspaceNumber),
          loadBoards(),
          workspaceApi.members(workspaceNumber),
        ]);
        if (cancelled) return;
        setCurrentWorkspace(workspace);
        setBoards(boardList);
        setMembers(memberList);
        const me = memberList.find((member: Member) => member.user_id === user?.id);
        setRole(me?.role ?? 'viewer');
        const active = boardList.find((board) => location.pathname.includes(`/boards/${board.id}`)) ?? boardList[0];
        if (active) setCurrentBoard(active);
      } catch (error) {
        message.error(errorMessage(error));
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
  }, [workspaceNumber, user?.id, location.pathname, setCurrentWorkspace, setCurrentBoard, setMembers]);

  if (loading && !currentWorkspace) {
    return (
      <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Spin size="large" />
      </div>
    );
  }

  const menuItems = [
    {
      key: 'boards',
      icon: <AppstoreOutlined />,
      label: '看板',
      children: boards.map((board) => ({
        key: `board-${board.id}`,
        label: board.name,
      })),
    },
    { key: 'stats', icon: <BarChartOutlined />, label: '统计' },
    { key: 'members', icon: <TeamOutlined />, label: '成员' },
  ];

  const selectedKey = location.pathname.includes('/stats')
    ? 'stats'
    : location.pathname.includes('/members')
      ? 'members'
      : currentBoard
        ? `board-${currentBoard.id}`
        : 'boards';

  const onMenuClick = ({ key }: { key: string }) => {
    if (key === 'stats') navigate(`/workspaces/${workspaceNumber}/stats`);
    else if (key === 'members') navigate(`/workspaces/${workspaceNumber}/members`);
    else if (key === 'boards') navigate(`/workspaces/${workspaceNumber}/boards`);
    else if (key.startsWith('board-')) {
      const board = boards.find((item) => `board-${item.id}` === key);
      if (board) {
        setCurrentBoard(board);
        navigate(`/workspaces/${workspaceNumber}/boards/${board.id}`);
      }
    }
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={240} theme="dark">
        <div style={{ padding: 16, color: '#fff', fontWeight: 700, fontSize: 18 }}>
          {currentWorkspace?.name}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={onMenuClick}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <Typography.Text strong>
            {currentBoard ? currentBoard.name : '团队任务看板'}
          </Typography.Text>
          <NotificationBell />
        </Header>
        <Content style={{ margin: 16 }}>
          <Outlet context={{ workspaceId: workspaceNumber, role, boards, setCurrentBoard, refreshBoards: loadBoards }} />
        </Content>
      </Layout>
    </Layout>
  );
}
