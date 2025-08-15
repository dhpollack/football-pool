import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import ProfilePage from './ProfilePage';

describe('ProfilePage', () => {
  const mockProfile = { name: 'John Doe', address: '123 Main St' };

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockProfile),
      })
    ) as vi.Mock;
  });

  it('renders the profile', async () => {
    render(<ProfilePage />);

    expect(await screen.findByDisplayValue(/john doe/i)).toBeInTheDocument();
    expect(await screen.findByDisplayValue(/123 main st/i)).toBeInTheDocument();
  });

  it('updates the profile', async () => {
    render(<ProfilePage />);

    await screen.findByDisplayValue(/john doe/i);

    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'Jane Doe' } });
    fireEvent.click(screen.getByText(/save/i));

    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/users/me',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ name: 'Jane Doe', address: '123 Main St' }),
      })
    );
  });
});