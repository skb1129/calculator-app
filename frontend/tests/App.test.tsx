import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import App from '../src/App';
import * as apiService from '../src/services/api';

describe('App Component', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders OmniCalc header, backend status indicator, and calculator', async () => {
    vi.spyOn(apiService, 'checkHealth').mockResolvedValue(true);

    render(<App />);

    expect(screen.getByText('OmniCalc')).toBeInTheDocument();
    expect(screen.getByText('Full-Stack')).toBeInTheDocument();
    expect(screen.getByLabelText('Calculator expression input')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText('Connected')).toBeInTheDocument();
    });
  });

  it('shows Disconnected status when backend health check fails', async () => {
    vi.spyOn(apiService, 'checkHealth').mockResolvedValue(false);

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('Disconnected')).toBeInTheDocument();
    });
  });
});
