import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { getUserStats } from '../services/user';
import { getUserGames } from '../services/game';
import './Dashboard.css';

const Dashboard = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState(null);
  const [recentGames, setRecentGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      if (!user?.id) {
        return;
      }

      try {
        setLoading(true);
        setError(null);

        // Fetch user statistics and recent games in parallel
        const [statsData, gamesData] = await Promise.all([
          getUserStats(user.id),
          getUserGames()
        ]);

        setStats(statsData);
        
        // Get the 5 most recent games
        if (gamesData.games) {
          const sortedGames = [...gamesData.games]
            .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
            .slice(0, 5);
          setRecentGames(sortedGames);
        }
      } catch (err) {
        console.error('Error fetching dashboard data:', err);
        setError('Failed to load dashboard data. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [user]);

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) {
      return 'Today';
    } else if (days === 1) {
      return 'Yesterday';
    } else if (days < 7) {
      return `${days} days ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  const getResultBadgeClass = (result) => {
    if (result === 'win') return 'badge-win';
    if (result === 'loss') return 'badge-loss';
    return 'badge-draw';
  };

  const getResultText = (result) => {
    if (result === 'win') return 'Won';
    if (result === 'loss') return 'Lost';
    return 'Draw';
  };

  if (loading) {
    return (
      <div className="dashboard-container">
        <div className="loading-state">
          <div className="spinner"></div>
          <p>Loading your dashboard...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="dashboard-container">
        <div className="error-state">
          <p>{error}</p>
          <button onClick={() => window.location.reload()} className="btn btn-primary">
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <div className="dashboard-content">
        <div className="dashboard-header">
          <h1>Welcome back, {user?.username}!</h1>
          <p className="subtitle">Ready for a challenge?</p>
        </div>

        {/* New Game CTA */}
        <div className="new-game-section">
          <button 
            onClick={() => navigate('/game')}
            className="btn btn-primary btn-large"
          >
            Start New Game
          </button>
        </div>

        {/* Statistics Cards */}
        <div className="stats-grid">
          <div className="stat-card card">
            <div className="stat-icon">🎮</div>
            <div className="stat-value">{stats?.total_games || 0}</div>
            <div className="stat-label">Total Games</div>
          </div>

          <div className="stat-card card">
            <div className="stat-icon">🏆</div>
            <div className="stat-value">{stats?.wins || 0}</div>
            <div className="stat-label">Wins</div>
          </div>

          <div className="stat-card card">
            <div className="stat-icon">❌</div>
            <div className="stat-value">{stats?.losses || 0}</div>
            <div className="stat-label">Losses</div>
          </div>

          <div className="stat-card card">
            <div className="stat-icon">🤝</div>
            <div className="stat-value">{stats?.draws || 0}</div>
            <div className="stat-label">Draws</div>
          </div>

          <div className="stat-card card highlight">
            <div className="stat-icon">📊</div>
            <div className="stat-value">
              {stats?.win_rate ? `${stats.win_rate.toFixed(1)}%` : '0%'}
            </div>
            <div className="stat-label">Win Rate</div>
          </div>
        </div>

        {/* Recent Games */}
        <div className="recent-games-section">
          <h2>Recent Games</h2>
          {recentGames.length === 0 ? (
            <div className="card empty-state">
              <p>No games yet. Start your first game above!</p>
            </div>
          ) : (
            <div className="games-list">
              {recentGames.map((game) => (
                <div key={game.id} className="game-item card">
                  <div className="game-info">
                    <span className={`result-badge ${getResultBadgeClass(game.result || game.status)}`}>
                      {getResultText(game.result || game.status)}
                    </span>
                    <span className="difficulty-badge">
                      {game.difficulty}
                    </span>
                  </div>
                  <div className="game-date">
                    {formatDate(game.created_at)}
                  </div>
                </div>
              ))}
            </div>
          )}
          {recentGames.length > 0 && (
            <div className="view-all">
              <Link to="/profile">View all games →</Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
