window.initCaptcha = function() {
  const container = document.getElementById("g-recaptcha");

  if (!container || container.innerHTML !== "") return;

  if (window.grecaptcha && window.grecaptcha.enterprise) {
    grecaptcha.enterprise.ready(function() {
      try {
        grecaptcha.enterprise.render(container, {
          sitekey: container.getAttribute("data-sitekey"),
        });
      } catch (e) {
        console.error("Captcha render failed:", e);
      }
    });
  }
};

document.getElementById("copy-year").textContent =
  new Date().getFullYear();

function renderMath(root) {
  if (typeof window.renderMathInElement === "function") {
    window.renderMathInElement(root || document.body, {
      delimiters: [
        { left: "$$", right: "$$", display: true },
        { left: "\\[", right: "\\]", display: true },
        { left: "\\(", right: "\\)", display: false },
      ],
      throwOnError: false,
    });
  }
}

document.addEventListener("DOMContentLoaded", function() {
  renderMath(document.body);
});

document.body.addEventListener("htmx:afterSwap", function(evt) {
  // On navigation, the script is usually already loaded.
  // We check if 'g-recaptcha' exists. If it does, we render.
  // If it DOESN'T exist, we do nothing—the 'onload' callback from Step 2 will handle it when it arrives.
  if (window.grecaptcha) {
    window.initCaptcha();
  }

  const target = evt.detail && evt.detail.target ? evt.detail.target : document.body;

  // Re-run Prism code highlighting on newly swapped content
  if (window.Prism) {
    window.Prism.highlightAllUnder(target);
  }

  // Re-run KaTeX math rendering on newly swapped content
  renderMath(target);
});

function getPreferredTheme() {
  try {
    const saved = localStorage.getItem("theme");
    if (saved === "dark" || saved === "light") {
      return saved;
    }
  } catch (e) {}

  if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
    return "dark";
  }
  return "light";
}

function applyTheme(theme) {
  const t = theme === "dark" ? "dark" : "light";
  document.documentElement.setAttribute("data-theme", t);
  if (document.body) {
    document.body.setAttribute("data-theme", t);
  }
  try {
    localStorage.setItem("theme", t);
  } catch (e) {}

  const isDark = t === "dark";
  const label = isDark ? "Switch to light mode" : "Switch to dark mode";
  document.querySelectorAll("[data-theme-toggle]").forEach(function(btn) {
    btn.setAttribute("aria-label", label);
    btn.setAttribute("title", label);
    btn.setAttribute("aria-pressed", isDark ? "true" : "false");
  });
}

// Event delegation handles clicks on toggle button or its child SVGs
document.addEventListener("click", function(event) {
  const toggleBtn = event.target.closest("[data-theme-toggle]");
  if (!toggleBtn) return;
  event.preventDefault();
  const current = document.documentElement.getAttribute("data-theme") || getPreferredTheme();
  const next = current === "dark" ? "light" : "dark";
  applyTheme(next);
});

// Follow system theme changes if no explicit user preference
if (window.matchMedia) {
  window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", function(e) {
    try {
      if (localStorage.getItem("theme")) return;
    } catch (err) {}
    applyTheme(e.matches ? "dark" : "light");
  });
}

// Initial application on script load
applyTheme(getPreferredTheme());

/**
* Highlight the nav link for the current page.
* Called on load and after every htmx navigation / history restore.
*/
function updateActiveNav() {
  const path = window.location.pathname;
  const section = path === "/" ? "home" : path.split("/")[1];

  document.querySelectorAll("[data-nav]").forEach((link) => {
    if (link.dataset.nav === section) {
      link.setAttribute("aria-current", "page");
    } else {
      link.removeAttribute("aria-current");
    }
  });
}

updateActiveNav();
document.body.addEventListener("htmx:pushedIntoHistory", updateActiveNav);
document.body.addEventListener("htmx:historyRestore", function() {
  updateActiveNav();
  if (window.Prism) {
    window.Prism.highlightAllUnder(document);
  }
  renderMath(document.body);
});
window.addEventListener("popstate", updateActiveNav);
