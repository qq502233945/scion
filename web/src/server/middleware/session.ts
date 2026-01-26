/**
 * Session middleware
 *
 * Configures koa-session for session management with secure cookie settings
 */

import type Koa from 'koa';
import session from 'koa-session';

import type { AppConfig } from '../config.js';
import type { User } from '../../shared/types.js';

/**
 * Session data stored in the session
 */
export interface SessionData {
  /** Authenticated user */
  user?: User;
  /** OAuth return URL after login */
  returnTo?: string;
  /** OAuth state for CSRF protection */
  oauthState?: string;
}

/**
 * Augment Koa's session type with our custom data
 */
declare module 'koa-session' {
  interface Session extends SessionData {}
}

/**
 * Session configuration options
 */
export interface SessionConfig {
  /** Session key (cookie name) */
  key: string;
  /** Max age in milliseconds */
  maxAge: number;
  /** Whether to use secure cookies (HTTPS only) */
  secure: boolean;
  /** HTTP only cookies (not accessible via JavaScript) */
  httpOnly: boolean;
  /** SameSite attribute */
  sameSite: 'strict' | 'lax' | 'none';
  /** Whether cookies are signed */
  signed: boolean;
}

/**
 * Get session configuration from app config
 */
export function getSessionConfig(config: AppConfig): SessionConfig {
  return {
    key: 'scion:sess',
    maxAge: config.session.maxAge,
    secure: config.production,
    httpOnly: true,
    sameSite: 'lax',
    signed: true,
  };
}

/**
 * Create session middleware
 *
 * @param app - Koa application instance
 * @param config - Application configuration
 * @returns Session middleware
 */
export function createSessionMiddleware(app: Koa, config: AppConfig): Koa.Middleware {
  // Validate session secret in production
  if (config.production && !config.session.secret) {
    throw new Error('SESSION_SECRET must be set in production');
  }

  // Set app keys for signed cookies
  app.keys = [config.session.secret];

  const sessionConfig = getSessionConfig(config);

  // Create session middleware with koa-session options
  return session(
    {
      key: sessionConfig.key,
      maxAge: sessionConfig.maxAge,
      httpOnly: sessionConfig.httpOnly,
      signed: sessionConfig.signed,
      secure: sessionConfig.secure,
      sameSite: sessionConfig.sameSite,
      // Disable auto commit to allow manual session management
      autoCommit: true,
      // Renew session if more than half the maxAge has passed
      renew: true,
    },
    app
  );
}
