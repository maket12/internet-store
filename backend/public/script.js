const apiUrl = "http://localhost:8080/api/v1";

document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("login-form").addEventListener("submit", login);
  document.getElementById("register-form").addEventListener("submit", register);
  document.getElementById("logout-button").addEventListener("click", logout);
  document
    .getElementById("add-product-form")
    .addEventListener("submit", addProduct);
  fetchProducts();
  fetchUserInfo();
});

async function login(event) {
  event.preventDefault();
  const username = document.getElementById("login-username").value;
  const password = document.getElementById("login-password").value;

  const response = await fetch(`${apiUrl}/login`, {
    credentials: "include",
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (response.ok) {
    alert("работает!");
    fetchProducts();
    fetchUserInfo();
  } else {
    const error = await response.json();
    alert(error.error);
  }
}

async function register(event) {
  event.preventDefault();
  const username = document.getElementById("register-username").value;
  const password = document.getElementById("register-password").value;

  const response = await fetch(`${apiUrl}/register`, {
    credentials: "include",
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (response.ok) {
    alert("Зарегался. Теперь логинься.");
  } else {
    const error = await response.json();
    alert(error.error);
  }
}

async function logout() {
  const response = await fetch(`${apiUrl}/logout`, {
    method: "POST",
    credentials: "include",
  });

  if (response.ok) {
    document.getElementById("user-info").classList.add("hidden");
    document.getElementById("logout-button").classList.add("hidden");
    document.getElementById("add-product-section").classList.add("hidden");
  } else {
    const error = await response.json();
    alert(error.error);
  }
}

async function fetchProducts() {
  const response = await fetch(`${apiUrl}/get_products`);
  const products = await response.json();
  displayProducts(products);
}

function displayProducts(products) {
  const tableBody = document.getElementById("product-table-body");
  tableBody.innerHTML = "";

  products.forEach((product) => {
    const row = document.createElement("tr");
    row.innerHTML = `
    <td style="text-align: center;">${product.id}</td>
    <td contenteditable="true" style="text-align: center;">${product.name}</td>
    <td contenteditable="true" style="text-align: center;">${product.description}</td>
    <td contenteditable="true" style="text-align: center;">${product.price}</td>
    <td contenteditable="true" style="text-align: center;">${product.available}</td>
    <td style="text-align: center;">
      <button onclick="deleteProduct('${product.id}')" class="bg-red-500 text-white rounded text-xs py-1 px-2">Удалить</button>
      <button onclick="saveProduct('${product.id}', this)" class="bg-blue-500 text-white rounded text-xs py-1 px-2">Сохранить</button>
    </td>
    `;
    tableBody.appendChild(row);
  });
}

async function deleteProduct(productId) {
  const response = await fetch(`${apiUrl}/remove_product`, {
    credentials: "include",
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id: productId }),
  });

  if (response.ok) {
    fetchProducts();
  } else {
    const error = await response.json();
    alert(error.error);
  }
}

async function saveProduct(productId, button) {
  const row = button.closest("tr");
  const newName = row.cells[1].innerText;
  const newDescription = row.cells[2].innerText;
  const newPrice = row.cells[3].innerText;
  const newAvailable = row.cells[4].innerText;

  const response = await fetch(`${apiUrl}/update_product`, {
    credentials: "include",
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      id: productId,
      name: newName,
      description: newDescription,
      price: parseFloat(newPrice),
      available: parseInt(newAvailable),
    }),
  });

  if (response.ok) {
    fetchProducts();
    document.getElementById("save-message").classList.remove("hidden");
    setTimeout(() => {
      document.getElementById("save-message").classList.add("hidden");
    }, 3000);
  } else {
    const error = await response.json();
    alert(error.error);
  }
}

async function fetchUserInfo() {
  const response = await fetch(`${apiUrl}/user_info`, {
    method: "POST",
    credentials: "include",
  });

  if (response.ok) {
    const user = await response.json();
    document.getElementById(
      "user-info-content"
    ).innerText = `Залогинился как: ${user.username} (Роль: ${user.role})`;
    document.getElementById("user-info").classList.remove("hidden");
    document.getElementById("logout-button").classList.remove("hidden");
    if (user.role === "admin") {
      document.getElementById("add-product-section").classList.remove("hidden");
    }
  }
}

async function addProduct(event) {
  event.preventDefault();
  const name = document.getElementById("product-name").value;
  const description = document.getElementById("product-description").value;
  const image = document.getElementById("product-image").value;
  const price = document.getElementById("product-price").value;
  const available = document.getElementById("product-available").value;

  const response = await fetch(`${apiUrl}/add_product`, {
    credentials: "include",
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name,
      description,
      image,
      price: parseFloat(price),
      available: parseInt(available),
    }),
  });

  if (response.ok) {
    fetchProducts();
  } else {
    const error = await response.json();
    alert(error.error);
  }
}
