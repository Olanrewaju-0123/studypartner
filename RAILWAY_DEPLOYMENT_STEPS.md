# Railway Deployment - Step by Step

## The Issue You Encountered

Railway was trying to build from the root directory, but your app has separate `backend/` and `frontend/` directories. You need to deploy them as separate services.

## Solution: Deploy as Separate Services

### Step 1: Access Your Railway Project

1. Go to [railway.app](https://railway.app)
2. You should see your "studypartner" project
3. Click on it to open the project dashboard

### Step 2: Add PostgreSQL Database

1. In your project dashboard, click **"Add Service"**
2. Select **"Database"** → **"PostgreSQL"**
3. Wait for the database to deploy (this may take a few minutes)
4. Once deployed, click on the PostgreSQL service
5. Go to the **"Variables"** tab
6. Copy the `DATABASE_URL` value (you'll need this for the backend)

### Step 3: Deploy Backend Service

1. In your project dashboard, click **"Add Service"**
2. Select **"GitHub Repo"**
3. Choose your repository
4. **IMPORTANT**: In the service settings, set the **Root Directory** to `backend`
5. Railway will detect it's a Go project and start building
6. Go to the **"Variables"** tab and add:
   ```
   DATABASE_URL=<paste the DATABASE_URL from PostgreSQL service>
   JWT_SECRET=your-secret-key-here
   HUGGINGFACE_API_KEY=your-huggingface-key
   PORT=8080
   ```
7. Wait for the backend to deploy

### Step 4: Deploy Frontend Service

1. In your project dashboard, click **"Add Service"**
2. Select **"GitHub Repo"**
3. Choose the same repository
4. **IMPORTANT**: In the service settings, set the **Root Directory** to `frontend`
5. Railway will detect it's a Node.js project and start building
6. Go to the **"Variables"** tab and add:
   ```
   NEXT_PUBLIC_API_URL=<backend-service-url>
   ```
   (You'll get the backend URL from the backend service's "Deployments" tab)

### Step 5: Configure Custom Domains

1. Go to each service's **"Settings"** → **"Domains"**
2. Click **"Generate Domain"** for both services
3. Update the frontend's `NEXT_PUBLIC_API_URL` to use the backend's custom domain

## Key Points to Remember

1. **Root Directory**: Always set the root directory to `backend` or `frontend` when adding services
2. **Environment Variables**: Make sure to set all required environment variables
3. **Database URL**: Use the DATABASE_URL from the PostgreSQL service
4. **API URL**: Update the frontend's API URL to point to the backend service

## Troubleshooting

- If a service fails to build, check the build logs in the "Deployments" tab
- Make sure all environment variables are set correctly
- Verify that the root directory is set correctly for each service

## Your Project Structure

```
studypartner/
├── backend/          ← Deploy as separate service (Root Directory: backend)
├── frontend/         ← Deploy as separate service (Root Directory: frontend)
├── docker-compose.yml
└── other files...
```

This approach will work because Railway will build each service from its respective directory.
