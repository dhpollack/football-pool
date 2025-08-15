
import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import PickEntryPage from './PickEntryPage';

describe('PickEntryPage', () => {
  const mockGames = [
    { id: 1, favorite_team: 'Team A', underdog_team: 'Team B', spread: 3 },
    { id: 2, favorite_team: 'Team C', underdog_team: 'Team D', spread: 7 },
  ];

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockGames),
      })
    ) as vi.Mock;
  });

  it('renders the games', async () => {
    render(<PickEntryPage />);

    expect(await screen.findByText(/team a/i)).toBeInTheDocument();
    expect(await screen.findByText(/team b/i)).toBeInTheDocument();
    expect(await screen.findByText(/team c/i)).toBeInTheDocument();
    expect(await screen.findByText(/team d/i)).toBeInTheDocument();
  });

  it('handles quick pick', async () => {
    render(<PickEntryPage />);

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/quick pick/i));

    // After quick pick, the selects should have values
    const selects = await screen.findAllByRole('combobox');
    selects.forEach((select) => {
      expect(select).not.toHaveValue('');
    });
  });

  it('submits picks', async () => {
    render(<PickEntryPage />);

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/submit picks/i));

    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/picks',
      expect.any(Object)
    );
  });
});
