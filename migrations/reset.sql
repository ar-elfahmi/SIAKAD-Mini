TRUNCATE TABLE
    enrollments,
    students,
    users,
    courses
RESTART IDENTITY
CASCADE;

-- psql -U postgres -d NAMA_DATABASE -f migrations/reset.sql


-- psql -U postgres -d siakad_mini -f migrations/005_seed_users.sql
-- psql -U postgres -d siakad_mini -f migrations/006_seed_students.sql
-- psql -U postgres -d siakad_mini -f migrations/007_seed_courses.sql
-- psql -U postgres -d siakad_mini -f migrations/008_seed_enrollments.sql