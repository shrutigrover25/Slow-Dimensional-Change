# 🔥 SCD Backend Database & API Testing Guide
## Database aur API Testing ke liye Complete Guide

---

## 🚀 Quick Start / Jaldi Shuru Karne ke liye

### 1. Database Setup (Database ka Setup)

```bash
# PostgreSQL install aur start karne ke liye
sudo apt install -y postgresql postgresql-contrib
sudo service postgresql start

# Database aur user setup
sudo -u postgres psql -c "CREATE DATABASE mercor;"
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'Shruti@25';"
```

### 2. Application Start (Application Chalane ke liye)

```bash
# Dependencies install karo
go mod tidy

# Application start karo
go run cmd/main.go
```

Application ab `http://localhost:8080` par chal raha hai! 🎉

---

## 🗄️ Database Operations (Database ke Operations)

### Database me Connect karne ke liye:
```bash
sudo -u postgres psql -d mercor
```

### Tables dekhne ke liye:
```sql
-- All tables show karo
\dt

-- Specific table ka structure dekho
\d jobs
\d timelogs  
\d payment_line_items
```

### Data Query karne ke liye:
```sql
-- Jobs table ka data
SELECT id, uid, version, status, rate, title 
FROM jobs 
ORDER BY id, version;

-- Timelogs ka data
SELECT id, uid, version, contractor_id, start_time, end_time 
FROM timelogs;

-- Payment line items ka data  
SELECT id, uid, version, contractor_id, amount, issued_at 
FROM payment_line_items;
```

### Manual Data Insert karne ke liye:
```sql
-- Naya job create karo
INSERT INTO jobs (id, uid, version, status, rate, title, company_id, contractor_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    gen_random_uuid(), 
    1,
    'active',
    50.0,
    'New Manual Job',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    NOW(),
    NOW()
);

-- Naya timelog create karo
INSERT INTO timelogs (id, uid, version, contractor_id, start_time, end_time, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    gen_random_uuid(),
    1,
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '1 hour',
    NOW(),
    NOW()
);

-- Naya payment create karo
INSERT INTO payment_line_items (id, uid, version, contractor_id, amount, issued_at, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    gen_random_uuid(),
    1,
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    75.50,
    NOW(),
    NOW(),
    NOW()
);
```

---

## 🔧 API Testing (API ka Testing)

### Testing Script Use karo:
```bash
# Complete test run karo
./test_api.sh

# Specific tests
./test_api.sh jobs      # Job operations test
./test_api.sh timelogs  # Timelog operations test  
./test_api.sh payments  # Payment operations test
./test_api.sh scd       # SCD versioning test
./test_api.sh db        # Database operations guide
```

### Manual API Testing:

#### 🏢 Job Operations:

**Get Job by UID:**
```bash
curl -s http://localhost:8080/jobs/00000000-0000-0000-0000-000000000003 | jq '.'
```

**Create New Job:**
```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New Test Job",
    "status": "active", 
    "rate": 45.0,
    "companyId": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    "contractorId": "cccccccc-cccc-cccc-cccc-cccccccccccc"
  }' | jq '.'
```

**Update Job (Creates New Version):**
```bash
# Pehle job UID get karo, phir update karo
JOB_UID="your-job-uid-here"

curl -X PUT http://localhost:8080/jobs/$JOB_UID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Job Title",
    "status": "extended",
    "rate": 55.0
  }' | jq '.'
```

**Update Job Status Only:**
```bash
curl -X PUT "http://localhost:8080/jobs/$JOB_UID/status?status=extended" | jq '.'
```

#### ⏰ Timelog Operations:

**Get Timelog:**
```bash
curl -s http://localhost:8080/timelogs/1c2e2ca7-a69d-421b-b278-f7f83a49e7e5 | jq '.'
```

**Create New Timelog:**
```bash
curl -X POST http://localhost:8080/timelogs \
  -H "Content-Type: application/json" \
  -d '{
    "startTime": "2025-07-28T14:00:00Z",
    "endTime": "2025-07-28T16:00:00Z", 
    "contractorId": "cccccccc-cccc-cccc-cccc-cccccccccccc"
  }' | jq '.'
```

#### 💰 Payment Operations:

**Get Payment:**
```bash
curl -s http://localhost:8080/payment-line-items/de1dbf39-3e6c-4d3b-af19-4447e2c26571 | jq '.'
```

**Create New Payment:**
```bash
curl -X POST http://localhost:8080/payment-line-items \
  -H "Content-Type: application/json" \
  -d '{
    "contractorId": "cccccccc-cccc-cccc-cccc-cccccccccccc",
    "amount": 125.50,
    "issuedAt": "2025-07-28T15:00:00Z"
  }' | jq '.'
```

---

## 🔄 SCD (Slowly Changing Dimensions) samjhiye

### Key Concepts:
- **ID**: Logical entity ID (same across versions) - Ek hi entity ke sabhi versions
- **UID**: Unique version identifier (different for each version) - Har version ka alag UID  
- **Version**: Version number (increments with updates) - Update ke saath badhta hai

### Example:
```bash
# Same job ke different versions dekho
curl -s http://localhost:8080/jobs/00000000-0000-0000-0000-000000000001 | jq '{Version, Status, Rate}'
curl -s http://localhost:8080/jobs/00000000-0000-0000-0000-000000000002 | jq '{Version, Status, Rate}'  
curl -s http://localhost:8080/jobs/00000000-0000-0000-0000-000000000003 | jq '{Version, Status, Rate}'
```

### Database me SCD dekhiye:
```sql
-- Same ID ke saare versions dekho
SELECT id, uid, version, status, rate, created_at 
FROM jobs 
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
ORDER BY version;
```

---

## 🧪 Testing Scenarios (Testing ke scenarios)

### 1. Create-Read-Update-Delete (CRUD) Testing:
```bash
# 1. Create job
NEW_JOB=$(curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Job", "status": "active", "rate": 40.0, "companyId": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", "contractorId": "cccccccc-cccc-cccc-cccc-cccccccccccc"}')

# 2. Extract UID
JOB_UID=$(echo $NEW_JOB | jq -r '.UID')

# 3. Read job
curl -s http://localhost:8080/jobs/$JOB_UID | jq '.'

# 4. Update job (creates new version)
curl -X PUT http://localhost:8080/jobs/$JOB_UID \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Test Job", "status": "extended", "rate": 50.0}' | jq '.'
```

### 2. Versioning Testing:
```bash
# Original job dekho
curl -s http://localhost:8080/jobs/$JOB_UID | jq '{Version, Status, Rate}'

# Update karo - naya version banega
UPDATED_JOB=$(curl -X PUT http://localhost:8080/jobs/$JOB_UID \
  -H "Content-Type: application/json" \
  -d '{"title": "Second Update", "status": "active", "rate": 60.0}')

# Naya UID get karo
NEW_UID=$(echo $UPDATED_JOB | jq -r '.UID')

# Dono versions compare karo
echo "Original Version:"
curl -s http://localhost:8080/jobs/$JOB_UID | jq '{Version, Status, Rate, UID}'

echo "Updated Version:"  
curl -s http://localhost:8080/jobs/$NEW_UID | jq '{Version, Status, Rate, UID}'
```

### 3. Relationship Testing:
```bash
# Job ke saath timelog create karo
curl -X POST http://localhost:8080/timelogs \
  -H "Content-Type: application/json" \
  -d "{
    \"startTime\": \"$(date -u -d '2 hours ago' '+%Y-%m-%dT%H:%M:%SZ')\",
    \"endTime\": \"$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')\",
    \"contractorId\": \"cccccccc-cccc-cccc-cccc-cccccccccccc\"
  }"
```

---

## 🚨 Troubleshooting (Problem Solving)

### Common Issues:

**1. Database connection issue:**
```bash
# PostgreSQL running hai check karo
sudo service postgresql status

# Agar nahi chal raha to start karo
sudo service postgresql start
```

**2. API not responding:**
```bash
# Application running hai check karo
ps aux | grep main

# Port check karo
netstat -tlnp | grep 8080

# Application restart karo
go run cmd/main.go
```

**3. Database access denied:**
```bash
# Password reset karo
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'Shruti@25';"
```

**4. Seed data duplicate error:**
- Yeh normal hai agar application multiple times start kiya ho
- Tables clear karne ke liye: `DROP TABLE jobs, timelogs, payment_line_items CASCADE;`

---

## 📝 Quick Commands Reference

```bash
# Database commands
sudo -u postgres psql -d mercor
\dt                              # Show tables
\d jobs                          # Show job table structure  
SELECT * FROM jobs LIMIT 5;     # Show sample data
\q                               # Quit

# API commands  
curl http://localhost:8080/jobs/[UID]                    # GET job
curl -X POST http://localhost:8080/jobs -d '{...}'       # CREATE job
curl -X PUT http://localhost:8080/jobs/[UID] -d '{...}'  # UPDATE job

# Testing commands
./test_api.sh                    # Run all tests
./test_api.sh --help             # Show help
./test_api.sh scd                # Test versioning
```

---

## 🎯 Next Steps (Aage ke Steps)

1. **Advanced Testing**: Integration tests likhiye
2. **Performance**: Load testing karo large datasets ke saath  
3. **Security**: Authentication aur authorization add karo
4. **Monitoring**: Logging aur metrics setup karo
5. **Documentation**: API documentation generate karo (Swagger)

---

## 📞 Support

Agar koi problem aaye to:
1. Application logs check karo: `tail -f /var/lib/postgresql/17/main/postgresql.log`
2. API status check karo: `curl -s http://localhost:8080/jobs/00000000-0000-0000-0000-000000000003`
3. Database connectivity test karo: `sudo -u postgres psql -d mercor -c "SELECT version();"`

Happy Testing! 🚀