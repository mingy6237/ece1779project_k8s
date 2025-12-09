# User Guide - Inventory Management System

## Overview

This document describes all pages in the inventory management frontend, their functions, and how to use them.

---

## Pages Overview

### 1. Root Page (`/`)

**Location**: [`app/page.tsx`](../app/page.tsx)

**Function**: Entry point that redirects users to the dashboard.

**Usage**:

- Automatically redirects to `/dashboard` in demo mode
- In production, would redirect to `/login` if not authenticated

---

### 2. Login Page (`/login`)

**Location**: [`app/login/page.tsx`](../app/login/page.tsx)

**Function**: User authentication entry point.

**Features**:

- Email and password input fields
- "Remember Me" checkbox for persistent sessions
- Error handling for invalid credentials
- JWT token storage in localStorage

**Usage**:

1. Enter your email address
2. Enter your password
3. (Optional) Check "Remember Me" to stay logged in
4. Click "Sign In"
5. On success, redirects to dashboard
6. On failure, shows error message

**Demo Credentials**:

- Manager: `manager@example.com` / `password`
- Staff: `staff@example.com` / `password`

---

### 3. Dashboard Page (`/dashboard`)

**Location**: [`app/(protected)/dashboard/page.tsx`](../app/(protected)/dashboard/page.tsx)

**Function**: Main overview page showing key metrics and alerts.

**Features**:

- **Metrics Cards**:
  - Total SKUs in inventory
  - Low stock items count
  - Active users count
  - Total locations count
- **Low Stock Alerts Widget** (Manager only):
  - Lists items below threshold
  - Shows current quantity and threshold
  - Links to item detail pages
- **Quick Actions**:
  - "Add New Item" button → `/dashboard/items/new`
  - "Browse Inventory" button → `/dashboard/items`

**Usage**:

1. View high-level metrics at a glance
2. Check low-stock alerts (if Manager role)
3. Click on alert items to view details
4. Use quick action buttons to navigate

**Access Control**:

- All authenticated users can view metrics
- Only Managers see low-stock alerts

---

### 4. Items List Page (`/dashboard/items`)

**Location**: [`app/(protected)/dashboard/items/page.tsx`](../app/(protected)/dashboard/items/page.tsx)

**Function**: Browse, search, and filter all inventory items.

**Features**:

- **Data Table**:
  - Columns: SKU, Name, Category, Total Quantity, Locations, Status
  - Sortable columns (click header to sort)
  - Pagination controls (page size: 20, 50, 100)
  - Real-time row highlighting when updated
- **Search Bar**:
  - Debounced search (300ms delay)
  - Searches by SKU or name
- **Filters**:
  - Category dropdown (electronics, furniture, office-supplies, hardware, software)
  - Location dropdown (Warehouse A, Warehouse B, Store 1, Store 2)
  - Quantity range (min/max)
- **Actions**:
  - "Add New Item" button → `/dashboard/items/new`
  - "Export CSV" button (downloads current view)
  - "View" button per row → `/dashboard/items/[id]`
- **Real-Time Updates**:
  - Green highlight animation when item updated
  - Connection status indicator in header

**Usage**:

1. **Browse Items**: Scroll through paginated list
2. **Search**: Type in search bar to find by SKU/name
3. **Filter**:
   - Select category to filter by type
   - Select location to see items at specific location
   - Set min/max quantity to find stock levels
4. **Sort**: Click column headers to sort ascending/descending
5. **Export**: Click "Export CSV" to download filtered results
6. **View Details**: Click eye icon to open item detail page
7. **Change Page Size**: Select 20, 50, or 100 items per page
8. **Navigate Pages**: Use pagination controls at bottom

**Access Control**: All authenticated users

---

### 5. Item Detail Page (`/dashboard/items/[id]`)

**Location**: [`app/(protected)/dashboard/items/[id]/page.tsx`](../app/(protected)/dashboard/items/[id]/page.tsx)

**Function**: View detailed information about a specific item.

**Features**:

- **Tabs Navigation**:
  1. Overview
  2. Locations
  3. Stock History
  4. Audit Log

- **Overview Tab**:
  - Item image (if available)
  - SKU, Name, Category
  - Description
  - Total quantity across all locations
  - Low stock threshold
  - Status badge (In Stock / Low Stock / Out of Stock)
  - Action buttons:
    - "Edit" → `/dashboard/items/[id]/edit`
    - "Adjust Stock" (opens modal)
    - "Delete Item" (with confirmation)

- **Locations Tab**:
  - List of all locations with quantities
  - Visual quantity indicators
  - Location addresses

- **Stock History Tab**:
  - Time-series chart showing stock levels over time
  - Period selectors: 7d, 30d, 90d, 1y
  - Interactive chart (hover for details)

- **Audit Log Tab**:
  - Timeline of all stock adjustments
  - Shows: user, location, delta, reason, timestamp
  - Chronological order (newest first)

**Usage**:

1. **View Overview**: Check item details and current stock
2. **Edit Item**: Click "Edit" to modify item details
3. **Adjust Stock**:
   - Click "Adjust Stock"
   - Select location
   - Choose Add or Remove
   - Enter quantity
   - Provide reason
   - Submit
4. **Delete Item**:
   - Click "Delete"
   - Confirm deletion (irreversible)
5. **Check Locations**: Switch to Locations tab to see distribution
6. **View Trends**: Switch to Stock History tab, select time period
7. **Review Changes**: Switch to Audit Log tab for complete history

**Real-Time Updates**: Item quantities update live when changed by other users

**Access Control**:

- All users can view
- Only Managers can edit or delete

---

### 6. Create Item Page (`/dashboard/items/new`)

**Location**: [`app/(protected)/dashboard/items/new/page.tsx`](../app/(protected)/dashboard/items/new/page.tsx)

**Function**: Add a new item to inventory.

**Features**:

- **Form Sections**:
  1. Basic Information
  2. Image (optional)
  3. Initial Stock Levels
  4. Location Quantities

- **Fields**:
  - SKU (required, unique, max 50 chars)
  - Name (required, max 200 chars)
  - Description (optional, multi-line)
  - Category (required, dropdown)
  - Threshold (optional, low-stock alert level)
  - Image URL (optional, with preview)
  - Initial quantities per location (0 or more)

- **Validation**:
  - Real-time field validation
  - Error messages shown below fields
  - Form disabled until valid

**Usage**:

1. **Enter Basic Info**:
   - SKU (e.g., "ELEC-001")
   - Name (e.g., "Wireless Mouse")
   - Description (optional)
2. **Select Category**: Choose from dropdown
3. **Set Threshold** (optional): Low-stock alert level
4. **Add Image** (optional):
   - Paste image URL
   - Preview shows automatically
   - Click X to remove
5. **Distribute Stock**:
   - Enter quantities for each location
   - Total calculated automatically
   - Can leave some at 0
6. **Submit**: Click "Create Item"
7. **Cancel**: Click "Cancel" to go back without saving

**Success**: Redirects to new item's detail page

**Access Control**: Managers only

---

### 7. Edit Item Page (`/dashboard/items/[id]/edit`)

**Location**: [`app/(protected)/dashboard/items/[id]/edit/page.tsx`](../app/(protected)/dashboard/items/[id]/edit/page.tsx)

**Function**: Modify existing item details.

**Features**:

- Pre-filled form with current item data
- Same fields as create page (except SKU cannot be changed)
- Editable:
  - Name
  - Description
  - Category
  - Threshold
  - Image URL
- **Note**: Quantities are adjusted via "Adjust Stock" modal, not here

**Usage**:

1. Form loads with current values
2. Modify desired fields
3. Click "Save Changes" to update
4. Click "Cancel" to discard changes
5. Validation errors shown inline

**Access Control**: Managers only

---

### 8. Audit Log Page (`/dashboard/audit`)

**Location**: [`app/(protected)/dashboard/audit/page.tsx`](../app/(protected)/dashboard/audit/page.tsx)

**Function**: View complete history of all inventory adjustments.

**Features**:

- **Filters**:
  - Item search (by SKU or name, client-side)
  - Location dropdown
  - Start date
  - End date
- **Timeline View**:
  - Chronological list of adjustments
  - Each entry shows:
    - Item name and SKU
    - User who made change
    - Location
    - Delta (e.g., +10, -5)
    - Old quantity → New quantity
    - Reason
    - Reference (if provided)
    - Timestamp
- **Pagination**:
  - Cursor-based (20 per page)
  - Previous/Next buttons
  - Shows "X entries" count

**Usage**:

1. **Browse All Logs**: Scroll through timeline
2. **Filter by Item**: Type in search box
3. **Filter by Location**: Select from dropdown
4. **Filter by Date Range**:
   - Set start date
   - Set end date
   - Click "Apply Filters"
5. **Clear Filters**: Click "Clear Filters" to reset
6. **Navigate Pages**: Use Previous/Next buttons

**Access Control**: All authenticated users

---

### 9. Alerts Page (`/dashboard/alerts`)

**Location**: [`app/(protected)/dashboard/alerts/page.tsx`](../app/(protected)/dashboard/alerts/page.tsx)

**Function**: Manage low-stock alerts and notifications.

**Features**:

- **Tabs**:
  - All
  - Open
  - Acknowledged
  - Resolved
- **Alert Cards**:
  - Item name and SKU
  - Severity badge (Warning / Critical)
  - Status badge (Open / Acknowledged / Resolved)
  - Current quantity
  - Threshold
  - Location
  - Created timestamp
- **Actions per Alert**:
  - "Acknowledge" (marks as seen)
  - "Resolve" (marks as fixed)
  - "Snooze" (hide for 24 hours)
- **Severity Levels**:
  - **Warning**: Stock below threshold but > 0
  - **Critical**: Stock = 0 (out of stock)

**Usage**:

1. **View Alerts**: See all active alerts
2. **Filter by Status**: Click tabs (All/Open/Acknowledged/Resolved)
3. **Acknowledge Alert**:
   - Click "Acknowledge" to mark as seen
   - Status changes to "Acknowledged"
4. **Resolve Alert**:
   - Click "Resolve" when stock replenished
   - Status changes to "Resolved"
5. **Snooze Alert**:
   - Click "Snooze" to hide for 24 hours
   - Alert reappears after snooze period

**Real-Time Updates**: New alerts appear automatically via WebSocket

**Access Control**:

- Managers: Full access
- Staff: View only (no actions)

---

## Common UI Components

### Navigation Menu

**Location**: [`app/(protected)/layout.tsx`](../app/(protected)/layout.tsx)

**Menu Items**:

1. Dashboard (home icon)
2. Items (package icon)
3. Audit Log (clipboard icon)
4. Alerts (bell icon) - with badge showing open count

**User Menu** (top right):

- User name and role
- Logout button

**Connection Status**:

- Green dot: WebSocket connected
- Yellow dot: Disconnected

---

### Inventory Adjust Modal

**Location**: [`components/inventory-adjust-modal.tsx`](../components/inventory-adjust-modal.tsx)

**Function**: Adjust stock quantity for an item at a specific location.

**Features**:

- Mode toggle: Add / Remove
- Location selector
- Quantity input (number)
- Reason input (required)
- Reference input (optional, e.g., PO number)
- Current quantity display
- New quantity preview

**Usage**:

1. Opened from item detail page or items list
2. Select location
3. Choose Add or Remove mode
4. Enter quantity
5. Enter reason (required)
6. Enter reference (optional)
7. Click "Submit"
8. Modal closes on success
9. Item updates immediately (optimistic UI)

**Validation**:

- Cannot remove more than current quantity
- Reason is required
- Quantity must be > 0

---

## Keyboard Shortcuts

- `Tab` / `Shift+Tab`: Navigate form fields
- `Enter`: Submit forms (when focused on input)
- `Esc`: Close modals/dialogs
- `/`: Focus search bar (items page)

---

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

---

## Mobile Support

All pages are responsive and work on mobile devices:

- Touch-friendly buttons and inputs
- Collapsible navigation menu
- Stacked layouts on small screens
- Optimized table views

---

## Accessibility Features

- ARIA labels on all interactive elements
- Keyboard navigation support
- Focus indicators
- High contrast mode compatible
- Screen reader friendly
- Error announcements

---

## Tips & Best Practices

1. **Search Efficiently**: Use specific SKUs for exact matches
2. **Filter Before Export**: Apply filters to export only needed data
3. **Check Audit Logs**: Review history before making large adjustments
4. **Set Thresholds**: Configure low-stock thresholds to get timely alerts
5. **Acknowledge Alerts**: Keep alert list clean by acknowledging seen items
6. **Use References**: Add PO numbers or order IDs when adjusting stock
7. **Monitor Real-Time**: Watch for green highlights indicating changes

---

## Troubleshooting

### Page Won't Load

- Check internet connection
- Verify backend API is running
- Clear browser cache

### Can't See Updates

- Check WebSocket connection indicator (should be green)
- Refresh page
- Verify backend WebSocket server is running

### Form Won't Submit

- Check for validation errors (red text below fields)
- Ensure all required fields are filled
- Verify unique constraints (e.g., duplicate SKU)

### Export Not Working

- Check browser popup blocker
- Ensure you have data to export (not filtered to 0 results)

---

## Contact Support

For issues not covered in this guide, contact:

- Frontend Team: Yuan Wang, Yanrong Xiao
- Backend Team: Yukun Zhou
- DevOps: Ming Yang
