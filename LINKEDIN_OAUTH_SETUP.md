# LinkedIn OAuth Setup Guide

## Overview

This guide walks you through setting up LinkedIn OAuth for social login in your application. Users can now sign in with their LinkedIn account, and we automatically collect their profile information including:

- First name
- Last name
- Profile picture URL
- Primary email address
- Email verification status
- Language/locale

## Prerequisites

1. A LinkedIn Developer Account
2. Access to LinkedIn App Console (https://www.linkedin.com/developers/apps)
3. Your application running locally (localhost:3000 for frontend, localhost:8080 for backend)

## Step 1: Create LinkedIn App

1. Go to [LinkedIn Developers Dashboard](https://www.linkedin.com/developers/apps)
2. Click **Create app**
3. Fill in the required information:
   - **App name**: "Recession Recovery Workspace" (or your app name)
   - **LinkedIn Page**: Select or create a LinkedIn company page
   - **Legal agreement**: Accept and create app
4. Once created, go to the **Auth** tab

## Step 2: Configure Authorized Redirect URLs

In the LinkedIn App Settings:

1. Find **"Authorized redirect URLs for your app"** section
2. Add these URLs:
   ```
   http://localhost:3000/auth/linkedin/callback
   http://localhost:8080/auth/linkedin/callback
   ```
3. For production, add your actual domain:
   ```
   https://yourdomain.com/auth/linkedin/callback
   ```

## Step 3: Get Your Credentials

In the **Auth** tab, find:

- **Client ID**: Copy this value
- **Client secret**: Copy this value (keep it secret!)

## Step 4: Configure Environment Variables

### Backend Configuration (backend/.env)

```env
LINKEDIN_CLIENT_ID=your_client_id_here
LINKEDIN_CLIENT_SECRET=your_client_secret_here
LINKEDIN_REDIRECT_URI=http://localhost:3000/auth/linkedin/callback
```

### Frontend Configuration (frontend/.env.local)

```env
NEXT_PUBLIC_LINKEDIN_CLIENT_ID=your_client_id_here
NEXT_PUBLIC_LINKEDIN_REDIRECT_URI=http://localhost:3000/auth/linkedin/callback
```

**Important**: The frontend only needs the CLIENT_ID and REDIRECT_URI. The CLIENT_SECRET must never be exposed to the frontend.

## Step 5: Request Access to OpenID Connect

To collect user profile data, you need to request access to OpenID Connect:

1. In LinkedIn App Console, go to the **Products** tab
2. Look for "OpenID Connect" in the products list
3. Click **Request access**
4. You should see "Sign In with LinkedIn" is available
5. In the **Auth** tab, enable the following scopes:
   - `openid` ✓
   - `profile` ✓
   - `email` ✓

## Step 6: Test the OAuth Flow

### Login with LinkedIn

1. Start your backend: `make backend-run`
2. Start your frontend: `make frontend-run`
3. Navigate to http://localhost:3000
4. On the login page, click "Login with LinkedIn"
5. You'll be redirected to LinkedIn
6. Approve the permissions request
7. You'll be redirected back to your app and automatically logged in

## Data Collected from LinkedIn

When a user authenticates via LinkedIn, we automatically collect and store:

```json
{
  "id": "user_id_from_linkedin",
  "sub": "linkedin_unique_identifier",
  "email": "user@example.com",
  "email_verified": true,
  "given_name": "John",
  "family_name": "Doe",
  "picture": "https://media.licdn.com/...",
  "locale": "en_US"
}
```

This data is stored in the database and linked to the user account.

## Database Changes

The User model has been extended with these new fields:

```go
type User struct {
    ID                    string // User ID (unchanged)
    Name                  string // Full name (unchanged)
    Email                 string // Email (unchanged)
    Role                  string // User role (unchanged)
    FirstName             string // LinkedIn first name
    LastName              string // LinkedIn last name
    ProfilePictureURL     string // LinkedIn profile picture URL
    LinkedInID            string // LinkedIn unique ID
    EmailVerified         bool   // Email verification status
    Locale                string // User's locale (language/region)
    AuthProvider          string // "email" or "linkedin"
}
```

## API Endpoints

### Get LinkedIn Auth URL
```
GET /api/v1/auth/linkedin/url?state={random_state}

Response:
{
  "url": "https://www.linkedin.com/oauth/v2/authorization?..."
}
```

### LinkedIn Callback Handler
```
POST /api/v1/auth/linkedin/callback

Request Body:
{
  "code": "authorization_code_from_linkedin",
  "state": "state_parameter_for_csrf_protection"
}

Response:
{
  "token": "jwt_token",
  "user": {
    "id": "user_id",
    "email": "user@example.com",
    "name": "John Doe",
    "first_name": "John",
    "last_name": "Doe",
    "profile_picture_url": "...",
    "email_verified": true,
    "auth_provider": "linkedin"
  }
}
```

## Security Considerations

1. **CSRF Protection**: We use the `state` parameter to prevent CSRF attacks
2. **Secret Management**: Never commit `.env` files with sensitive credentials
3. **HTTPS**: Always use HTTPS in production
4. **Token Validation**: JWT tokens are validated on every authenticated request
5. **Scope Limitation**: We only request the minimum necessary scopes

## Troubleshooting

### "Redirect URI Mismatch"
- Ensure the redirect URI in LinkedIn App Console exactly matches your `LINKEDIN_REDIRECT_URI`
- No trailing slashes unless you included them

### "Invalid authorization code"
- The code may have expired (valid for 10 minutes)
- User may have denied permissions
- Try the flow again

### "Profile picture not loading"
- LinkedIn requires your app to be in production mode for profile picture access
- For development, profile pictures may not be available

### User not getting created
- Check that email verification is working
- Verify database migrations have run: `make migrate`
- Check backend logs for errors

## Next Steps

1. ✅ Configure LinkedIn App
2. ✅ Set environment variables
3. ✅ Test the login flow
4. ✅ Customize user profile page to show LinkedIn data
5. ✅ Deploy to production

## Production Deployment

For production:

1. Update LinkedIn App's redirect URIs to your production domain
2. Update environment variables in production:
   - Backend: Set via deployment platform or environment
   - Frontend: Set via `NEXT_PUBLIC_` variables or build-time configuration
3. Use HTTPS for all URLs
4. Consider enabling LinkedIn's advanced security features

## Support

For issues with LinkedIn OAuth:
- Check [LinkedIn API Documentation](https://docs.microsoft.com/en-us/linkedin/)
- Review [OpenID Connect Documentation](https://openid.net/connect/)
- Check application logs in `logs/backend.log`

## Switching Between Email and LinkedIn Login

Users can choose either method:
- **Email/Password**: Traditional registration and login
- **LinkedIn**: One-click social login

If a user tries to sign up with LinkedIn using an email that already exists in the system (from email signup), the accounts are automatically linked.
