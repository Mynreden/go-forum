document.querySelector(".profil").addEventListener("click", function () {
  var arrow = this.querySelector(".arrow-4");
  arrow.classList.toggle("open");

  var navBar = document.querySelector(".profil_nav_bar");
  navBar.classList.toggle("visible");
});

var prevScrollPos = window.pageYOffset;

window.onscroll = function () {
  var currentScrollPos = window.pageYOffset;

  if (prevScrollPos > currentScrollPos) {
    prevScrollPos++;
    // Scrolling up, show the navbar
    document.querySelector(".nav_bar").style.top = "63px";
  } else {
    prevScrollPos++;
    // Scrolling down, hide the navbar
    document.querySelector(".nav_bar").style.top = "-123px"; // Adjust this value based on the height of your navbar
  }
};

document.addEventListener("DOMContentLoaded", function () {
  createObserver();

  var textarea = document.getElementById("myTextarea");

  if (!textarea) {
    return;
  }

  textarea.addEventListener("input", function () {
    this.style.height = "auto";
    this.style.height = this.scrollHeight + "px";
  });
});

function toggleForm() {
  var createPostButton = document.getElementById("createPost");
  var postForm = document.getElementById("postForm");
  var textarea = document.getElementById("myTextarea");

  // Toggle visibility of elements
  createPostButton.style.display = "none";
  postForm.style.display = "block";

  // Set focus on textarea
  textarea.focus();
}

function adjustTextarea() {
  var textarea = document.getElementById("myTextarea");
  var submitPostBox = document.querySelector(".submit_post_box");
  var contentLength = Number(textarea?.value.length);

  // Decrease font size if character count is greater than 199
  textarea.style.fontSize = contentLength > 199 ? "15px" : "20px";

  // Dynamically set minimum height for .submit_post_box, considering textarea height and added padding
  submitPostBox.style.minHeight = textarea.scrollHeight + 20 + "px";
}

// Add input event listener for textarea
var textarea = document.getElementById("myTextarea");
if (textarea) {
  textarea.addEventListener("input", adjustTextarea);
}

// Additionally, respond to font size changes (e.g., outside input)
window.addEventListener("resize", adjustTextarea);

function previewImage(event) {
  var input = event.target;
  var preview = document.getElementById("image-preview");

  while (preview.firstChild) {
    preview.removeChild(preview.firstChild);
  }

  var files = input.files;
  if (files && files[0]) {
    var file = files[0];

    // Check file size (in bytes)
    var maxSize = 2 * 1024 * 1024; // 2MB
    if (file.size > maxSize) {
      alert("File size exceeds 2MB. Please choose a smaller file.");
      input.value = ""; // Clear the input
      return;
    }

    // Check file type
    var allowedTypes = ["image/jpeg", "image/png", "image/gif"];
    if (allowedTypes.indexOf(file.type) === -1) {
      alert(
        "Invalid file type. Please choose a valid image file (JPEG, PNG, or GIF)."
      );
      input.value = ""; // Clear the input
      return;
    }

    var reader = new FileReader();

    reader.onload = function (e) {
      var img = document.createElement("img");
      img.src = e.target.result;
      preview.appendChild(img);
    };

    reader.readAsDataURL(file);

    // Show the container when an image is selected
    preview.style.display = "block";
  } else {
    // Hide the container when no image is selected
    preview.style.display = "none";
  }
}

function toggleList() {
  var listContainer = document.getElementById("listContainer");
  listContainer.style.display =
    listContainer.style.display === "none" ? "block" : "none";
}

function addItem() {
  var checkboxes = document.querySelectorAll(
    '#listContainer input[type="checkbox"]:checked'
  );
  var inputField = document.getElementById("inputField");
  var selectedItems = [];

  checkboxes.forEach(function (checkbox) {
    selectedItems.push(checkbox.value);
  });

  // Update the input content with selected items
  inputField.value = selectedItems.join(", ");

  // Close the list after adding items to the input
  document.getElementById("listContainer").style.display = "none";
}

const ROUTES = {
  "/": "posts",
  "/likedposts": "lp",
  "/postscat": "pc",
  "/myposts": "mp",
};

function fetchNextPostsPage({ offset, limit }) {
  const params = new URLSearchParams();
  params.set("offset", offset);
  params.set("limit", limit);

  const existingParams = new URLSearchParams(window.location.search).toString();

  const currentRoute = ROUTES[window.location.pathname];

  if (!currentRoute) {
    console.error("STAS DOBAV' ROUTE");
    return;
  }

  return fetch(
    `${currentRoute}?${params.toString()}${
      existingParams ? "&" + existingParams : ""
    }`,
    {
      headers: {
        "INFINITE-SCROLL": true,
      },
    }
  ).then((data) => data.text());
}

function createObserver() {
  const limit = 10;
  let offset = 10;

  const posts = document.querySelector(".posts");

  const callback = async ([entry], observer) => {
    if (entry.isIntersecting) {
      const nextPage = await fetchNextPostsPage({
        offset: offset,
        limit: limit,
      });

      if (nextPage.length <= 1) {
        return;
      }

      posts.innerHTML += nextPage;
      offset += limit;
    }
  };

  const options = {
    root: null,
    threshold: 0,
  };

  const observer = new IntersectionObserver(callback, options);
  const footer = document.querySelector(".footer");

  observer.observe(footer);
}

function validateForm() {
  var categoryInput = document.getElementById("inputField");
  var titleInput = document.querySelector(".pst-tt-upl");
  var contentTextarea = document.getElementById("myTextarea");

  // Check if any of the required fields is empty
  if (
    !categoryInput.value.trim() ||
    !titleInput.value.trim() ||
    !contentTextarea.value.trim()
  ) {
    alert("Please fill in all required fields: Category, Title, and Content.");
    return false;
  }

  // Additional validation logic can be added here if needed

  return true; // Form is valid, allow submission
}

// script.js
document.addEventListener("DOMContentLoaded", function () {
  var categoriesContainer = document.querySelector(".post-ctg");
  var preElements = categoriesContainer.querySelectorAll(".post-category-anme");

  // Ограничиваем количество элементов до 3
  var maxPreElements = 3;
  if (preElements.length > maxPreElements) {
    // Удаляем лишние элементы
    for (var i = maxPreElements; i < preElements.length; i++) {
      preElements[i].remove();
    }

    // Добавляем многоточие к последнему элементу
    var lastPreElement = preElements[maxPreElements - 1];
    lastPreElement.textContent += "...";
  }
});


const deletePost = async (id) => {
  try {
    const response = await fetch(`/post/delete?id=${id}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      throw new Error(`Failed to delete post: ${response.statusText}`);
    }

    window.location.href = "/myposts"
  } catch (error) {
    console.error('Error deleting post:', error);
  }
};
