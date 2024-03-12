## Quick Start

1. **Start Dependencies**:

   This command starts the PostgreSQL and Redis services using Docker:

   ```sh
   make dependencies_start
   ```
2. **Build Frontend**:

   Install npm dependencies and build the frontend:

   ```sh
   make front_build
   ```
3. **Database Migration**:

   Run migrations and seed the database with initial data:

   ```sh
   make migrate_up
   ```
4. **Insert Data**:

   Seed the database using initial data:

   ```sh
   make insert_data
   ```
5. **Run the Server**:

   Start the backend server using the Go command or through the configured VSCode launch task:

   ```sh
   make run
   ```
   Or alternatively, if you're using VSCode with a configured launch.json, you can start the server using the IDE's debug feature.