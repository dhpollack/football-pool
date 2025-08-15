
import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import OverallResultsPage from './OverallResultsPage';

describe('OverallResultsPage', () => {
  it('renders the loading state', () => {
    render(<OverallResultsPage />);

    expect(screen.getByText(/loading overall results.../i)).toBeInTheDocument();
  });

  it('renders the error state', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        json: () => Promise.resolve({ message: 'Failed to fetch overall results' }),
      })
    ) as vi.Mock;

    render(<OverallResultsPage />);

    expect(await screen.findByText(/failed to fetch overall results/i)).toBeInTheDocument();
  });

  it('renders the results', async () => {
    const mockResults = [
      { player_name: 'Player 1', score: 100 },
      { player_name: 'Player 2', score: 90 },
    ];

    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockResults),
      })
    ) as vi.Mock;

    render(<OverallResultsPage />);

    expect(await screen.findByText(/player 1/i)).toBeInTheDocument();
    expect(await screen.findByText(/100/i)).toBeInTheDocument();
    expect(await screen.findByText(/player 2/i)).toBeInTheDocument();
    expect(await screen.findByText(/90/i)).toBeInTheDocument();
  });
});
