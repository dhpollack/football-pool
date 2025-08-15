
import { render, screen } from '@testing-library/react';
import HomePage from './HomePage';

describe('HomePage', () => {
  it('renders the welcome message', () => {
    render(<HomePage />);

    expect(screen.getByText(/welcome to the football pool!/i)).toBeInTheDocument();
  });
});
