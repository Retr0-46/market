const productsElement = document.querySelector('#products');
const statusElement = document.querySelector('#status');
const limitElement = document.querySelector('#limit');
const dialog = document.querySelector('#product-dialog');
const detailElement = document.querySelector('#product-detail');

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

function renderProducts(products) {
  productsElement.replaceChildren();
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
    detailElement.replaceChildren(category, art(product.type, 'detail-art'), name, amount, stock);
  } catch (error) {
    if (dialog.open) detailElement.textContent = `Не удалось загрузить товар. ${error.message}`;
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
loadProducts();
