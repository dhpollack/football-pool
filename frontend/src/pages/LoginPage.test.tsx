
import { render, screen } from '@testing-library/react';
import LoginPage from './LoginPage';

describe('LoginPage', () => {
  it('renders the login form', () => {
    render(<LoginPage />);

    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
  });
});
