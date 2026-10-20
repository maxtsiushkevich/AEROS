package http

import (
	"fmt"
	"net/http"
)

func (s *Server) rbacAdminPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>RBAC Admin</title>
  <style>
    :root {
      --bg: #0f172a;
      --panel: #111827;
      --card: #1f2937;
      --line: #374151;
      --text: #e5e7eb;
      --muted: #9ca3af;
      --accent: #60a5fa;
      --success: #34d399;
      --danger: #f87171;
      --warning: #fbbf24;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Arial, sans-serif;
      background: var(--bg);
      color: var(--text);
    }
    .container {
      max-width: 1280px;
      margin: 32px auto;
      padding: 24px;
    }
    h1, h2 { margin: 0 0 16px; }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 16px;
      margin-bottom: 24px;
    }
    .card {
      background: var(--card);
      border: 1px solid var(--line);
      border-radius: 12px;
      padding: 16px;
      box-shadow: 0 8px 24px rgba(0,0,0,0.15);
    }
    form {
      display: grid;
      gap: 10px;
    }
    label {
      display: grid;
      gap: 6px;
      font-size: 13px;
      color: var(--muted);
    }
    input, button {
      font: inherit;
      border-radius: 8px;
      border: 1px solid var(--line);
      padding: 10px 12px;
    }
    input {
      background: rgba(15, 23, 42, 0.8);
      color: var(--text);
    }
    button {
      cursor: pointer;
      background: var(--accent);
      color: var(--bg);
      font-weight: 700;
      border: none;
    }
    button.danger {
      background: var(--danger);
      color: white;
    }
    button.secondary {
      background: var(--warning);
      color: #111827;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 10px;
      font-size: 14px;
    }
    th, td {
      border-bottom: 1px solid var(--line);
      padding: 10px 8px;
      text-align: left;
      vertical-align: top;
    }
    th {
      color: var(--muted);
      font-weight: 600;
    }
    .status {
      margin-top: 12px;
      min-height: 20px;
      color: var(--muted);
      font-size: 13px;
    }
    .row {
      display: flex;
      gap: 8px;
      align-items: center;
    }
    .inline {
      display: inline-flex;
      gap: 8px;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>RBAC Admin</h1>
    <div id="status" class="status"></div>

    <div class="grid">
      <section class="card">
        <h2>Role</h2>
        <form id="role-form">
          <label>name<input name="name" required /></label>
          <label>description<input name="description" /></label>
          <button type="submit">Create role</button>
        </form>
      </section>

      <section class="card">
        <h2>Action</h2>
        <form id="action-form">
          <label>name<input name="name" required /></label>
          <button type="submit">Create action</button>
        </form>
      </section>

      <section class="card">
        <h2>Resource</h2>
        <form id="resource-form">
          <label>name<input name="name" required /></label>
          <label>description<input name="description" /></label>
          <button type="submit">Create resource</button>
        </form>
      </section>

      <section class="card">
        <h2>Permission</h2>
        <form id="permission-form">
          <label>resource<input name="resource_name" required /></label>
          <label>action<input name="action_name" required /></label>
          <button type="submit">Create permission</button>
        </form>
      </section>

      <section class="card">
        <h2>Grant permission to role</h2>
        <form id="grant-form">
          <label>role<input name="role_name" required /></label>
          <label>resource<input name="resource_name" required /></label>
          <label>action<input name="action_name" required /></label>
          <button type="submit">Grant</button>
        </form>
      </section>

      <section class="card">
        <h2>Assign role to user</h2>
        <form id="assign-form">
          <label>user_id<input name="user_id" required placeholder="UUID" /></label>
          <label>role<input name="role_name" required /></label>
          <button type="submit">Assign</button>
        </form>
      </section>
    </div>

    <div class="card">
      <h2>Roles</h2>
      <table>
        <thead><tr><th>Name</th><th>Description</th><th>Action</th></tr></thead>
        <tbody id="roles-body"></tbody>
      </table>
    </div>

    <div class="grid" style="margin-top: 24px;">
      <div class="card">
        <h2>Actions</h2>
        <table><thead><tr><th>Name</th></tr></thead><tbody id="actions-body"></tbody></table>
      </div>
      <div class="card">
        <h2>Resources</h2>
        <table><thead><tr><th>Name</th><th>Description</th></tr></thead><tbody id="resources-body"></tbody></table>
      </div>
    </div>

    <div class="card" style="margin-top: 24px;">
      <h2>Permissions</h2>
      <table><thead><tr><th>Resource</th><th>Action</th></tr></thead><tbody id="permissions-body"></tbody></table>
    </div>
  </div>

  <script>
    const statusEl = document.getElementById('status');
    const apiBase = '/api/v1/rbac';

    function setStatus(message, isError = false) {
      statusEl.textContent = message;
      statusEl.style.color = isError ? '#f87171' : '#9ca3af';
    }

    function parseJSON(response) {
      if (response.status === 204) return null;
      return response.json();
    }

    async function request(path, options = {}) {
      const response = await fetch(apiBase + path, {
        headers: { 'Content-Type': 'application/json' },
        ...options
      });
      const data = await parseJSON(response);
      if (!response.ok) {
        const message = data && (data.error || data.message) ? (data.error || data.message) : 'HTTP ' + response.status;
        throw new Error(message);
      }
      return data;
    }

    function valueOr(obj, keys) {
      for (const key of keys) {
        if (obj && obj[key] !== undefined && obj[key] !== null && obj[key] !== '') {
          return obj[key];
        }
      }
      return '';
    }

    function roleRow(role) {
      const name = valueOr(role, ['name', 'Name']);
      const description = valueOr(role, ['description', 'Description']);
      return '<tr><td>' + name + '</td><td>' + description + '</td><td><button class="danger" data-remove-role="' + name + '">Delete</button></td></tr>';
    }

    function actionRow(action) {
      const name = valueOr(action, ['name', 'Name']);
      return '<tr><td>' + name + '</td></tr>';
    }

    function resourceRow(resource) {
      const name = valueOr(resource, ['name', 'Name']);
      const description = valueOr(resource, ['description', 'Description']);
      return '<tr><td>' + name + '</td><td>' + description + '</td></tr>';
    }

    function permissionRow(permission) {
      const resourceName = valueOr(permission, ['resource_name', 'ResourceName']);
      const actionName = valueOr(permission, ['action_name', 'ActionName']);
      return '<tr><td>' + resourceName + '</td><td>' + actionName + '</td></tr>';
    }

    async function loadAll() {
      try {
        const rolesPromise = request('/roles');
        const actionsPromise = request('/actions');
        const resourcesPromise = request('/resources');
        const permissionsPromise = request('/permissions');

        const [roles, actions, resources, permissions] = await Promise.all([
          rolesPromise,
          actionsPromise,
          resourcesPromise,
          permissionsPromise
        ]);

        document.getElementById('roles-body').innerHTML = (roles || []).map(roleRow).join('');
        document.getElementById('actions-body').innerHTML = (actions || []).map(actionRow).join('');
        document.getElementById('resources-body').innerHTML = (resources || []).map(resourceRow).join('');
        document.getElementById('permissions-body').innerHTML = (permissions || []).map(permissionRow).join('');
        setStatus('Data loaded');
      } catch (error) {
        setStatus(error.message, true);
      }
    }

    document.getElementById('role-form').addEventListener('submit', async function (event) {
      event.preventDefault();
      const form = new FormData(event.target);
      try {
        await request('/roles', {
          method: 'PUT',
          body: JSON.stringify({ name: form.get('name'), description: form.get('description') })
        });
        event.target.reset();
        await loadAll();
      } catch (error) {
        setStatus(error.message, true);
      }
    });

    document.getElementById('action-form').addEventListener('submit', async function (event) {
      event.preventDefault();
      const form = new FormData(event.target);
      try {
        await request('/actions', {
          method: 'PUT',
          body: JSON.stringify({ name: form.get('name') })
        });
        event.target.reset();
        await loadAll();
      } catch (error) {
        setStatus(error.message, true);
      }
    });

    document.getElementById('resource-form').addEventListener('submit', async function (event) {
      event.preventDefault();
      const form = new FormData(event.target);
      try {
        await request('/resources', {
          method: 'PUT',
          body: JSON.stringify({ name: form.get('name'), description: form.get('description') })
        });
        event.target.reset();
        await loadAll();
      } catch (error) {
        setStatus(error.message, true);
      }
    });

    document.getElementById('permission-form').addEventListener('submit', async function (event) {
      event.preventDefault();
      const form = new FormData(event.target);
      try {
        await request('/permissions', {
          method: 'PUT',
          body: JSON.stringify({ resource_name: form.get('resource_name'), action_name: form.get('action_name') })
        });
        event.target.reset();
        await loadAll();
      } catch (error) {
        setStatus(error.message, true);
      }
    });

    document.getElementById('grant-form').addEventListener('submit', async function (event) {
      event.preventDefault();
      const form = new FormData(event.target);
      const roleName = form.get('role_name');
      try {
        await request('/roles/' + encodeURIComponent(roleName) + '/permissions', {
          method: 'PUT',
          body: JSON.stringify({ resource_name: form.get('resource_name'), action_name: form.get('action_name') })
        });
        event.target.reset();
        await loadAll();
      } catch (error) {
        setStatus(error.message, true);
      }
    });

    document.getElementById('assign-form').addEventListener('submit', async function (event) {
      event.preventDefault();
      const form = new FormData(event.target);
      const userId = form.get('user_id');
      try {
        await request('/users/' + encodeURIComponent(userId) + '/roles', {
          method: 'PUT',
          body: JSON.stringify({ role_name: form.get('role_name') })
        });
        event.target.reset();
        await loadAll();
      } catch (error) {
        setStatus(error.message, true);
      }
    });

    document.addEventListener('click', async function (event) {
      const roleBtn = event.target.closest('[data-remove-role]');
      if (roleBtn) {
        const roleName = roleBtn.getAttribute('data-remove-role');
        try {
          await request('/roles/' + encodeURIComponent(roleName), { method: 'DELETE' });
          await loadAll();
        } catch (error) {
          setStatus(error.message, true);
        }
      }
    });

    loadAll();
  </script>
</body>
</html>`)
	}
}
