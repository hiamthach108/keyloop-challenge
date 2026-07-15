# 3-5 Minute Demo Script

## 1. Opening

Hi, this is my Scenario B submission for the Keyloop coding challenge.

The goal is to help a dealership manager monitor vehicle inventory, quickly
spot stock older than 90 days, and record follow-up actions such as price
reduction, transfer proposal, marketing campaign, review, or another custom
action.

## 2. Architecture Summary

I built this as a small full-stack modular monolith.

The frontend is a Next.js dashboard using Ant Design and TanStack Query. The
backend is a Go REST API using Echo, GORM repositories, PostgreSQL, and Atlas
migrations. Docker Compose runs everything together: Postgres, the migration
container, the backend, and the frontend.

I chose this shape because it is easy to run for a challenge, but still has
clear boundaries. HTTP handlers only handle requests. The inventory service
owns the business rules. Repositories own database access. PostgreSQL owns
constraints and durable history.

## 3. Data Model

The core data model separates vehicle identity from inventory stock.

Vehicles store stable information like VIN, make, model, and model year.
Inventory stock stores dealership, price, current status, stocked-in date, and
stocked-out date.

Then there are two history tables:

Stock movements record stock in and stock out events. Inventory actions record
manager decisions for aging stock.

This gives the dashboard a current state, but also keeps an audit trail.

## 4. Backend Behavior

The main API is:

`GET /api/v1/dealerships/{dealershipID}/stocks`

It supports search, age filters, status, pagination, and sorting.

The key business rule is simple: a stock item is aging when it is in stock and
older than 90 calendar days. The backend calculates this and returns `isAging`.

Actions can only be created for in-stock aging vehicles. If a vehicle is too
new or out of stock, the backend rejects the action.

For data setup, I use Atlas migrations. The schema is generated from models,
but migrations are real SQL files protected by `atlas.sum`. The seed data
creates 50 stock rows across two dealerships, so the UI has enough data for a
proper demo.

## 5. Frontend Walkthrough

In the dashboard, I can switch between dealerships.

The table shows the stock list and has one simple search box. I kept search as
a single field because for this challenge it is enough to test make and model
matching without overcomplicating the UI.

Pagination and total results are at the top of the table, so the table does not
create nested scrollbars.

On the right, I use Ant Design Splitter for the detail panel instead of a modal.
This keeps the selected row and the detail context visible at the same time.
The first row is selected automatically, so the panel is never empty after data
loads.

The detail panel shows stock information, current status, action form, and
history. When I save an action, TanStack Query invalidates the stock list and
history queries, so the table and detail panel refresh cleanly.

## 6. Codex Usage

I used Codex as a coding partner during the implementation.

First, I used it to turn the challenge brief into an architecture plan. Then I
iterated with it on the backend model when the first version was too simple. We
moved from a basic vehicle list to proper inventory stock, stock movements,
action history, enum action types, repository per model, and Atlas migrations.

On the frontend, I used Codex to connect the dashboard to the backend API,
replace modal actions with a Splitter detail panel, add TanStack Query,
simplify search, and clean up the table layout.

The important point is that Codex helped me move faster, but I still reviewed
the domain rules, data model, migrations, and UI behavior. I verified the final
work with Go tests, frontend lint/typecheck/build, Docker Compose, API smoke
tests, and browser checks.

## 7. Closing

The result is a practical dealership inventory dashboard:

- real backend persistence,
- repeatable migrations and seed data,
- dealership-scoped APIs,
- aging-stock business rules,
- action history,
- and a clean frontend workflow for reviewing and acting on stock.

That is the core of my Scenario B solution.
