# Railway Deployment Guide

## Railway Pricing & Free Tier

✅ **Good News**: Railway's free tier is **permanent** and doesn't expire!

- **$5 monthly credit** that renews every month
- **PostgreSQL database** included
- **Custom domains** available
- **Automatic deployments** from GitHub
- **No time limits** on the free tier

## Deployment Steps

### 1. Prepare Your Repository

Make sure your code is pushed to GitHub with all the recent changes:

```bash
git add .
git commit -m "Prepare for Railway deployment"
git push origin main
```

### 2. Create Railway Account

1. Go to [railway.app](https://railway.app)
2. Sign up with your GitHub account
3. Connect your GitHub repository

### 3. Deploy PostgreSQL Database

1. In Railway dashboard, click "New Project"
2. Select "Deploy from GitHub repo"
3. Choose your repository
4. Click "Add Service" → "Database" → "PostgreSQL"
5. Railway will automatically create a PostgreSQL database
6. Note the connection details (you'll need them for the backend)

### 4. Deploy Backend Service

1. In the same project, click "Add Service" → "GitHub Repo"
2. Select your repository
3. Railway will detect it's a Go project
4. Set the **Root Directory** to `backend`
5. Configure environment variables:
   ```
   DATABASE_URL=<from PostgreSQL service>
   JWT_SECRET=your-secret-key-here
   HUGGINGFACE_API_KEY=your-huggingface-key
   PORT=8080
   ```
6. Deploy the service

### 5. Deploy Frontend Service

1. Add another service → "GitHub Repo"
2. Select the same repository
3. Set the **Root Directory** to `frontend`
4. Configure environment variables:
   ```
   NEXT_PUBLIC_API_URL=<backend-service-url>
   ```
5. Deploy the service

### 6. Configure Custom Domains

1. Go to each service's "Settings" → "Domains"
2. Generate custom domains for both services
3. Update the frontend's `NEXT_PUBLIC_API_URL` to use the backend's custom domain

## Environment Variables

### Backend (.env)

```
DATABASE_URL=postgres://username:password@host:port/database
JWT_SECRET=your-secret-key-here
HUGGINGFACE_API_KEY=your-huggingface-key
PORT=8080
```

### Frontend (.env.local)

```
NEXT_PUBLIC_API_URL=https://your-backend-domain.railway.app
```

## Monitoring & Maintenance

- **Logs**: Check service logs in Railway dashboard
- **Metrics**: Monitor usage in the dashboard
- **Updates**: Push to GitHub to trigger automatic deployments
- **Health Checks**: Both services have health check endpoints

## Cost Management

- Monitor your usage in the Railway dashboard
- The $5 credit should be sufficient for a small to medium application
- If you exceed the credit, you'll be charged per usage
- You can set spending limits in your account settings

## Troubleshooting

1. **Build Failures**: Check the build logs in Railway dashboard
2. **Database Connection**: Verify DATABASE_URL is correct
3. **CORS Issues**: Ensure frontend URL is in backend CORS settings
4. **Environment Variables**: Double-check all required variables are set

## Zero Downtime Deployments

Railway supports zero-downtime deployments by default:

- New instances start before old ones stop
- Health checks ensure new instances are ready
- Rolling deployment strategy minimizes downtime

Your application will remain available during updates!
