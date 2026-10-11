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

document.body.addEventListener("htmx:afterSwap", function() {
  // On navigation, the script is usually already loaded.
  // We check if 'g-recaptcha' exists. If it does, we render.
  // If it DOESN'T exist, we do nothing—the 'onload' callback from Step 2 will handle it when it arrives.
  if (window.grecaptcha) {
    window.initCaptcha();
  }
});

/**
* Utility function to calculate the current theme setting.
* Look for a local storage value.
* Fall back to system setting.
* Fall back to light mode.
*/
function calculateSettingAsThemeString({ localStorageTheme, systemSettingDark }) {
  if (localStorageTheme !== null) {
    return localStorageTheme;
  }

  if (systemSettingDark.matches) {
    return "dark";
  }

  return "light";
}

/**
* Utility function to update the button text and aria-label.
*/
function updateButton({ buttonEl, isDark }) {
  const newCta = isDark ? true : false;
  buttonEl.checked = newCta;
}

/**
* Utility function to update the theme setting on the html tag
*/
function updateThemeOnHtmlEl({ theme }) {
  document.querySelector("html").setAttribute("data-theme", theme);
}


/**
* On page load:
*/

/**
* 1. Grab what we need from the DOM and system settings on page load
*/
const button = document.querySelector("[data-theme-toggle]");
const localStorageTheme = localStorage.getItem("theme");
const systemSettingDark = window.matchMedia("(prefers-color-scheme: dark)");

/**
* 2. Work out the current site settings
*/
let currentThemeSetting = calculateSettingAsThemeString({ localStorageTheme, systemSettingDark });

/**
* 3. Update the theme setting and button text accoridng to current settings
*/
updateButton({ buttonEl: button, isDark: currentThemeSetting === "dark" });
updateThemeOnHtmlEl({ theme: currentThemeSetting });

/**
* 4. Add an event listener to toggle the theme
*/
button.addEventListener("change", (event) => {
  const newTheme = currentThemeSetting === "dark" ? "light" : "dark";

  localStorage.setItem("theme", newTheme);
  updateButton({ buttonEl: button, isDark: newTheme === "dark" });
  updateThemeOnHtmlEl({ theme: newTheme });

  currentThemeSetting = newTheme;
});
