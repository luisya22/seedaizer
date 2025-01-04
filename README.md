# Seedaizer: AI-Powered Database Seeder
Seedaizer is an innovative tool designed to automate and simplify populating relational databases with realistic mock data. It utilizes AI to generate SQL INSERT statements that respect your database's schema and relationships.

## Features
- **AI-Powered Data Generation:** Generate realistic and diverse mock data that aligns with your schema.
- **Schema-Aware:** Automatically understands your database structure, including tables, columns, and relationships.
- **Customizable:** Define data generation rules and patterns to suit your needs.
- **Cross-Database Support:** Works with popular relational databases like MySQL, PostgreSQL, SQLite, and more.
- **Developer-Friendly:** CLI-based for easy integration into your workflow.

## Usage
1. **Scan Database Schema:** Extract table structures and relationships.
```bash
seedaizer scandb --config ./config.yaml
```

2. **Generate Data:** Create SQL `INSERT` statements tailored to your database schema.
```bash
seedaizer seedup -q "insert 5 {users} with" -o print
```

3. **Configure:** USe a `config.yaml` file for database connection and other settings.
```yaml
database:
  url: "root:password@tcp(localhost:3306)/example_db"
openaikey: "your_openai_api_key"
```

## Coming Soon / TODO
