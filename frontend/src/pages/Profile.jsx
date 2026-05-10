import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { getUserStats } from '../services/user';
import { getUserGames } from '../services/game';
import './Profile.css';

const Profile = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState(null);
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filter, setFilter] = useState('all'); // all, win, loss, draw
  const [sortBy, setSortBy] = useState('date'); // date, difficulty

  useEffect(() => {
    const fetchData = async () => {
      if (!user?.id) {
        return;
      }

      try {
        setLoading(true);
        setError(null);

        // Fetch user statistics and all games in parallel
        const [statsData, gamesData] = await Promise.all([
          getUserStats(user.id),
          getUserGames()
        ]);

        setStats(statsData);
        
        // Get all completed games
        if (gamesData.games) {
          const completedGames = gamesData.games.filter(g => g.result);
          setGames(completedGames);
        }
      } catch (err) {
        console.error('Error fetching profile data:', err);
        setError('Failed to load profile data. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [user]);

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatCreatedDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
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

  const getDifficultyColor = (difficulty) => {
    switch (difficulty) {
      case 'easy': return 'difficulty-easy';
      case 'medium': return 'difficulty-medium';
      case 'hard': return 'difficulty-hard';
      case 'impossible': return 'difficulty-impossible';
      default: return '';
    }
  };

  // Filter games
  const filteredGames = games.filter(game => {
    if (filter === 'all') return true;
    return game.result === filter;
  });

  // Sort games
  const sortedGames = [...filteredGames].sort((a, b) => {
    if (sortBy === 'date') {
      return new Date(b.created_at) - new Date(a.created_at);
    } else {
      // Sort by difficulty
      const difficultyOrder = { easy: 1, medium: 2, hard: 3, impossible: 4 };
      return (difficultyOrder[b.difficulty] || 0) - (difficultyOrder[a.difficulty] || 0);
    }
  });

  // Group games by date
  const groupedGames = sortedGames.reduce((groups, game) => {
    const date = new Date(game.created_at).toDateString();
    if (!groups[date]) {
      groups[date] = [];
    }
    groups[date].push(game);
    return groups;
  }, {});

  if (loading) {
    return (
      <div className="profile-container">
        <div className="loading-state">
          <div className="spinner"></div>
          <p>Loading your profile...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="profile-container">
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
    <div className="profile-container">
      <div className="profile-content">
        {/* Profile Header */}
        <div className="profile-header card">
          <div className="profile-avatar">
            {user?.username?.charAt(0).toUpperCase()}
          </div>
          <div className="profile-info">
            <h1>{user?.username}</h1>
            <p className="profile-email">{user?.email}</p>
            <p className="profile-joined">
              Member since {formatCreatedDate(user?.created_at || new Date())}
            </p>
          </div>
        </div>

        {/* Detailed Statistics */}
        <div className="stats-section">
          <h2>Statistics</h2>
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
        </div>

        {/* Game History */}
        <div className="game-history-section">
          <div className="section-header">
            <h2>Game History</h2>
            <div className="game-count">
              {filteredGames.length} {filteredGames.length === 1 ? 'game' : 'games'}
            </div>
          </div>

          {/* Filters and Sorting */}
          <div className="controls-bar">
            <div className="filter-group">
              <label>Filter:</label>
              <button 
                className={`filter-btn ${filter === 'all' ? 'active' : ''}`}
                onClick={() => setFilter('all')}
              >
                All
              </button>
              <button 
                className={`filter-btn ${filter === 'win' ? 'active' : ''}`}
                onClick={() => setFilter('win')}
              >
                Wins
              </button>
              <button 
                className={`filter-btn ${filter === 'loss' ? 'active' : ''}`}
                onClick={() => setFilter('loss')}
              >
                Losses
              </button>
              <button 
                className={`filter-btn ${filter === 'draw' ? 'active' : ''}`}
                onClick={() => setFilter('draw')}
              >
                Draws
              </button>
            </div>

            <div className="sort-group">
              <label>Sort by:</label>
              <select 
                value={sortBy} 
                onChange={(e) => setSortBy(e.target.value)}
                className="sort-select"
              >
                <option value="date">Date</option>
                <option value="difficulty">Difficulty</option>
              </select>
            </div>
          </div>

          {/* Games List */}
          {sortedGames.length === 0 ? (
            <div className="card empty-state">
              <p>
                {filter === 'all' 
                  ? 'No games yet. Start playing to see your history!'
                  : `No ${filter === 'draw' ? 'draws' : filter + 's'} yet.`
                }
              </p>
            </div>
          ) : sortBy === 'date' ? (
            // Grouped by date
            <div className="games-timeline">
              {Object.entries(groupedGames).map(([date, dateGames]) => (
                <div key={date} className="timeline-group">
                  <h3 className="timeline-date">{formatCreatedDate(dateGames[0].created_at)}</h3>
                  <div className="games-list">
                    {dateGames.map((game) => (
                      <div key={game.id} className="game-item card">
                        <div className="game-main">
                          <span className={`result-badge ${getResultBadgeClass(game.result)}`}>
                            {getResultText(game.result)}
                          </span>
                          <span className={`difficulty-badge ${getDifficultyColor(game.difficulty)}`}>
                            {game.difficulty.charAt(0).toUpperCase() + game.difficulty.slice(1)}
                          </span>
                        </div>
                        <div className="game-time">
                          {new Date(game.created_at).toLocaleTimeString('en-US', {
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            // Simple list when not sorted by date
            <div className="games-list">
              {sortedGames.map((game) => (
                <div key={game.id} className="game-item card">
                  <div className="game-main">
                    <span className={`result-badge ${getResultBadgeClass(game.result)}`}>
                      {getResultText(game.result)}
                    </span>
                    <span className={`difficulty-badge ${getDifficultyColor(game.difficulty)}`}>
                      {game.difficulty.charAt(0).toUpperCase() + game.difficulty.slice(1)}
                    </span>
                  </div>
                  <div className="game-date">
                    {formatDate(game.created_at)}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Profile;
