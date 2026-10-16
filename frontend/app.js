const productsElement = document.querySelector('#products');
const statusElement = document.querySelector('#status');
const limitElement = document.querySelector('#limit');
const dialog = document.querySelector('#product-dialog');
const detailElement = document.querySelector('#product-detail');
const authDialog = document.querySelector('#auth-dialog');
const authForm = document.querySelector('#auth-form');
const authError = document.querySelector('#auth-error');
const cartDialog = document.querySelector('#cart-dialog');
const cartItemsElement = document.querySelector('#cart-items');
const cartStatus = document.querySelector('#cart-status');
let authMode = 'login';
let accessToken = '';
let expiresAt = 0;
let currentUser = null;
let currentProducts = new Map();

const categories = {
  sport: { label: 'Спорт', icon: '⚽' },
  electronics: { label: 'Электроника', icon: '⌨️' },
  books: { label: 'Книги', icon: '📚' },
  clothes: { label: 'Одежда', icon: '👕' }
};

let selectedType = '';
let requestNumber = 0;
const formatter = new Intl.NumberFormat('ru-RU');
const moneyFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
  maximumFractionDigits: 2
});

function price(value) {
  return moneyFormatter.format(value / 100);
}

function art(type, extraClass = '') {
  const element = document.createElement('div');
  element.className = `product-art ${extraClass}`;
  element.dataset.type = type;
  element.textContent = categories[type]?.icon ?? '✳';
  return element;
}

function setSession(token, expiry) {
  accessToken = token;
  expiresAt = new Date(expiry).getTime();
  currentUser = null;
  document.querySelector('#user-name').hidden = true;
  document.querySelector('#auth-button').hidden = Boolean(token);
  document.querySelector('#logout-button').hidden = !token;
  document.querySelector('#cart-count').hidden = !token;
  document.querySelector('#cart-count').textContent = '0';
}

function hasSession() {
  if (accessToken && Date.now() < expiresAt) return true;
  if (accessToken) setSession('', 0);
  return false;
}

async function api(path, options = {}, authenticated = false) {
  const headers = { ...options.headers };
  if (options.body) headers['Content-Type'] = 'application/json';
  if (authenticated) {
    if (!hasSession()) throw new Error('Войдите в аккаунт');
    headers.Authorization = `Bearer ${accessToken}`;
  }
  const response = await fetch(path, { ...options, headers });
  if (response.status === 401 && authenticated) {
    setSession('', 0);
    throw new Error('Сессия истекла. Войдите снова.');
  }
  if (!response.ok) {
    const message = (await response.text()).trim();
    throw new Error(response.status >= 500 ? 'Сервис временно недоступен' : message || `Ошибка ${response.status}`);
  }
  return response.status === 204 ? null : response.json();
}

function openAuth(mode = 'login') {
  authMode = mode;
  authError.textContent = '';
  authForm.reset();
  document.querySelector('#auth-title').textContent = mode === 'register' ? 'Регистрация' : 'Вход';
  document.querySelector('#auth-submit').textContent = mode === 'register' ? 'Зарегистрироваться' : 'Войти';
  document.querySelector('#password').autocomplete = mode === 'register' ? 'new-password' : 'current-password';
  for (const [name, id] of [['login', 'login-tab'], ['register', 'register-tab']]) {
    const button = document.querySelector(`#${id}`);
    button.classList.toggle('active', mode === name);
    button.setAttribute('aria-pressed', String(mode === name));
  }
  if (!authDialog.open) authDialog.showModal();
}

async function loadUser() {
  try {
    const user = await api('/auth/me', {}, true);
    currentUser = user;
    const name = document.querySelector('#user-name');
    name.textContent = user.username;
    name.hidden = false;
  } catch (error) {
    if (authDialog.open) authError.textContent = error.message;
  }
}

async function loadCartCount() {
  if (!hasSession()) return;
  try {
    const data = await api('/cart', {}, true);
    document.querySelector('#cart-count').textContent = (data.items ?? []).reduce((sum, item) => sum + item.quantity, 0);
  } catch (error) {
    if (cartDialog.open) cartStatus.textContent = error.message;
  }
}

function renderProducts(products) {
  productsElement.replaceChildren();
  currentProducts = new Map(products.map(product => [product.id, product]));
  for (const product of products) {
    const card = document.createElement('button');
    card.className = 'product-card';
    card.type = 'button';

    const info = document.createElement('div');
    info.className = 'product-info';
    const category = document.createElement('p');
    category.className = 'product-category';
    category.textContent = categories[product.type]?.label ?? product.type;
    const name = document.createElement('p');
    name.className = 'product-name';
    name.textContent = product.name;
    const amount = document.createElement('p');
    amount.className = 'product-price';
    amount.textContent = price(product.price);
    info.append(category, name, amount);
    card.append(art(product.type), info);
    card.addEventListener('click', () => showProduct(product.id));
    productsElement.append(card);
  }
}

async function loadProducts() {
  const currentRequest = ++requestNumber;
  statusElement.textContent = 'Загружаем товары…';
  productsElement.replaceChildren();

  try {
    const response = await fetch('/products/list', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ number: Number(limitElement.value), p_types: selectedType ? [selectedType] : [] })
    });
    if (!response.ok) throw new Error(`Ошибка ${response.status}`);
    const data = await response.json();
    if (currentRequest !== requestNumber) return;
    const products = Array.isArray(data.product) ? data.product : [];
    renderProducts(products);
    statusElement.textContent = products.length ? '' : 'Пока нет товаров в этой категории.';
  } catch (error) {
    if (currentRequest === requestNumber) statusElement.textContent = `Не удалось загрузить товары. ${error.message}`;
  }
}

async function showProduct(id) {
  detailElement.replaceChildren();
  detailElement.textContent = 'Загружаем товар…';
  dialog.showModal();

  try {
    const response = await fetch(`/products?id=${encodeURIComponent(id)}`);
    if (!response.ok) throw new Error(`Ошибка ${response.status}`);
    const product = await response.json();
    if (!dialog.open) return;

    const category = document.createElement('p');
    category.className = 'product-category';
    category.textContent = categories[product.type]?.label ?? product.type;
    const name = document.createElement('h3');
    name.className = 'detail-name';
    name.id = 'dialog-title';
    name.textContent = product.name;
    const amount = document.createElement('p');
    amount.className = 'detail-price';
    amount.textContent = price(product.price);
    const stock = document.createElement('p');
    stock.className = 'detail-stock';
    stock.textContent = `В наличии: ${formatter.format(product.stock_quantity)} шт.`;
    const actions = document.createElement('div');
    actions.className = 'detail-actions';
    const addButton = document.createElement('button');
    addButton.className = 'primary-button';
    addButton.type = 'button';
    addButton.textContent = 'Добавить в корзину';
    addButton.disabled = product.stock_quantity <= 0;
    addButton.addEventListener('click', async () => {
      if (!hasSession()) {
        dialog.close();
        openAuth('login');
        return;
      }
      addButton.disabled = true;
      try {
        await api('/cart/items', {
          method: 'POST',
          body: JSON.stringify({ product_id: product.id, quantity: 1 })
        }, true);
        addButton.textContent = 'Добавлено';
        await loadCartCount();
      } catch (error) {
        addButton.textContent = error.message;
      } finally {
        addButton.disabled = false;
      }
    });
    actions.append(addButton);
    detailElement.replaceChildren(category, art(product.type, 'detail-art'), name, amount, stock, actions);
  } catch (error) {
    if (dialog.open) detailElement.textContent = `Не удалось загрузить товар. ${error.message}`;
  }
}

async function showCart() {
  if (!hasSession()) {
    openAuth('login');
    return;
  }
  if (!cartDialog.open) cartDialog.showModal();
  await loadCart();
}

async function loadCart() {
  if (!cartDialog.open || !hasSession()) return;
  cartStatus.textContent = 'Загружаем корзину…';
  cartItemsElement.replaceChildren();
  document.querySelector('#cart-footer').hidden = true;
  try {
    const data = await api('/cart', {}, true);
    const items = data.items ?? [];
    if (!items.length) {
      cartStatus.textContent = 'Корзина пуста.';
      await loadCartCount();
      return;
    }
    cartStatus.textContent = '';
    let total = 0;
    for (const item of items) {
      let product = currentProducts.get(item.product_id);
      if (!product) {
        try {
          product = await api(`/products?id=${encodeURIComponent(item.product_id)}`);
        } catch { /* Товар может больше не существовать. */ }
      }
      const row = document.createElement('div');
      row.className = 'cart-item';
      const main = document.createElement('div');
      main.className = 'cart-item-main';
      const name = document.createElement('p');
      name.className = 'cart-item-name';
      name.textContent = product?.name ?? `Товар #${item.product_id}`;
      const amount = document.createElement('p');
      amount.className = 'cart-item-price';
      amount.textContent = product ? `${price(product.price)} × ${item.quantity}` : `${item.quantity} шт.`;
      if (product) total += product.price * item.quantity;
      main.append(name, amount);

      const actions = document.createElement('div');
      actions.className = 'cart-item-actions';
      const quantity = document.createElement('input');
      quantity.className = 'quantity-input';
      quantity.type = 'number';
      quantity.min = '1';
      quantity.step = '1';
      quantity.value = item.quantity;
      quantity.setAttribute('aria-label', `Количество: ${name.textContent}`);
      quantity.addEventListener('change', async () => {
        const value = Number(quantity.value);
        if (!Number.isSafeInteger(value) || value < 1) {
          quantity.value = item.quantity;
          cartStatus.textContent = 'Количество должно быть положительным числом.';
          return;
        }
        quantity.disabled = true;
        try {
          await api(`/cart/items/${item.product_id}`, { method: 'PUT', body: JSON.stringify({ quantity: value }) }, true);
          await loadCart();
        } catch (error) {
          quantity.disabled = false;
          quantity.value = item.quantity;
          cartStatus.textContent = error.message;
        }
      });
      const remove = document.createElement('button');
      remove.className = 'text-button';
      remove.type = 'button';
      remove.textContent = 'Удалить';
      remove.addEventListener('click', async () => {
        remove.disabled = true;
        try {
          await api(`/cart/items/${item.product_id}`, { method: 'DELETE' }, true);
          await loadCart();
        } catch (error) {
          remove.disabled = false;
          cartStatus.textContent = error.message;
        }
      });
      actions.append(quantity, remove);
      row.append(main, actions);
      cartItemsElement.append(row);
    }
    document.querySelector('#cart-total').textContent = `Итого: ${price(total)}`;
    document.querySelector('#cart-footer').hidden = false;
    await loadCartCount();
  } catch (error) {
    cartStatus.textContent = error.message;
  }
}

document.querySelectorAll('.filter').forEach(button => {
  button.addEventListener('click', () => {
    document.querySelectorAll('.filter').forEach(item => {
      const active = item === button;
      item.classList.toggle('active', active);
      item.setAttribute('aria-pressed', String(active));
    });
    selectedType = button.dataset.type;
    loadProducts();
  });
});

limitElement.addEventListener('change', loadProducts);
document.querySelector('#close-dialog').addEventListener('click', () => dialog.close());
dialog.addEventListener('click', event => { if (event.target === dialog) dialog.close(); });
document.querySelector('#auth-button').addEventListener('click', () => openAuth('login'));
document.querySelector('#logout-button').addEventListener('click', () => {
  setSession('', 0);
  cartDialog.close();
});
document.querySelector('#cart-button').addEventListener('click', showCart);
document.querySelector('#close-auth').addEventListener('click', () => authDialog.close());
document.querySelector('#close-cart').addEventListener('click', () => cartDialog.close());
authDialog.addEventListener('click', event => { if (event.target === authDialog) authDialog.close(); });
cartDialog.addEventListener('click', event => { if (event.target === cartDialog) cartDialog.close(); });
document.querySelector('#login-tab').addEventListener('click', () => openAuth('login'));
document.querySelector('#register-tab').addEventListener('click', () => openAuth('register'));
authForm.addEventListener('submit', async event => {
  event.preventDefault();
  authError.textContent = '';
  const submit = document.querySelector('#auth-submit');
  submit.disabled = true;
  const username = document.querySelector('#username').value.trim();
  const password = document.querySelector('#password').value;
  try {
    if (authMode === 'register') {
      await api('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ username, password, account_type: 'user' })
      });
    }
    const login = await api('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    });
    setSession(login.access_token, login.expires_at);
    authDialog.close();
    await Promise.all([loadUser(), loadCartCount()]);
  } catch (error) {
    authError.textContent = error.message;
  } finally {
    submit.disabled = false;
  }
});
document.querySelector('#clear-cart').addEventListener('click', async () => {
  const button = document.querySelector('#clear-cart');
  button.disabled = true;
  try {
    await api('/cart', { method: 'DELETE' }, true);
    await loadCart();
  } catch (error) {
    cartStatus.textContent = error.message;
  } finally {
    button.disabled = false;
  }
});
loadProducts();
