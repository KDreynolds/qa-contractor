(function () {
	"use strict";

	const html = document.documentElement;
	const toggle = document.getElementById("theme-toggle");

	/* ---- Theme ---- */
	function getTheme() {
		return localStorage.getItem("theme") || "dark";
	}

	function setTheme(t) {
		html.setAttribute("data-theme", t);
		localStorage.setItem("theme", t);
		if (toggle) toggle.innerHTML = t === "dark" ? "&#9788;" : "&#9790;";
	}

	setTheme(getTheme());

	if (toggle) {
		toggle.addEventListener("click", function () {
			setTheme(getTheme() === "dark" ? "light" : "dark");
		});
	}

	/* ---- Mobile nav toggle ---- */
	var navToggle = document.getElementById("nav-toggle");
	var navLinks = document.getElementById("nav-links");

	if (navToggle && navLinks) {
		navToggle.addEventListener("click", function () {
			navLinks.classList.toggle("open");
		});

		document.addEventListener("click", function (e) {
			if (
				!navToggle.contains(e.target) &&
				!navLinks.contains(e.target)
			) {
				navLinks.classList.remove("open");
			}
		});
	}

	/* ---- Contact form ---- */
	var contactForm = document.getElementById("contact-form");
	if (contactForm) {
		var fsId = contactForm.dataset.formspree;
		var btn = contactForm.querySelector("button[type=submit]");
		var btnText = document.getElementById("form-button-text");
		var indicator = document.getElementById("form-indicator");

		function showLoading() {
			if (btn) btn.disabled = true;
			if (btnText) btnText.style.display = "none";
			if (indicator) indicator.style.display = "";
		}

		function resetButton() {
			if (btn) btn.disabled = false;
			if (btnText) btnText.style.display = "";
			if (indicator) indicator.style.display = "none";
		}

		function showSuccess() {
			contactForm.outerHTML =
				'<div id="contact-form" class="form-success">' +
				'<div class="check">&#10003;</div>' +
				"<h3>Message sent!</h3>" +
				"<p>Thank you for reaching out. I&rsquo;ll get back to you within 24 hours.</p>" +
				"</div>";
		}

		contactForm.addEventListener("submit", function (e) {
			e.preventDefault();
			showLoading();

			if (fsId) {
				fetch("https://formspree.io/f/" + fsId, {
					method: "POST",
					body: new FormData(contactForm),
					headers: { Accept: "application/json" },
				})
					.then(function (resp) {
						if (resp.ok) {
							showSuccess();
						} else {
							resetButton();
							alert("Something went wrong. Please try again or email me directly.");
						}
					})
					.catch(function () {
						resetButton();
						alert("Network error. Please try again or email me directly.");
					});
			} else if (window.location.hostname === "localhost" || window.location.hostname === "127.0.0.1") {
				fetch("/contact", {
					method: "POST",
					body: new FormData(contactForm),
				})
					.then(function (resp) {
						if (resp.ok) return resp.text();
						throw new Error("Server error");
					})
					.then(function (html) {
						contactForm.outerHTML = html;
					})
					.catch(function () {
						resetButton();
						alert("Something went wrong. Please try again.");
					});
			} else {
				resetButton();
				alert("Contact form not configured. Please email me directly at kylereynoldsdev@gmail.com.");
			}
		});
	}

	/* ---- Active nav link on HTMX navigation ---- */
	function updateActiveNav() {
		var links = document.querySelectorAll(".nav-links a");
		var path = (window.location.pathname || "/").replace(/\/+$/, "") || "/";
		links.forEach(function (link) {
			var href = (link.getAttribute("href") || "").replace(/\/+$/, "") || "/";
			if (href === path) {
				link.classList.add("active");
			} else {
				link.classList.remove("active");
			}
		});
	}

	/* ---- HTMX lifecycles ---- */
	document.body.addEventListener("htmx:afterSettle", function () {
		updateActiveNav();
	});
	// Also handle popstate (browser back/forward) since hx-push-url is used
	window.addEventListener("popstate", function () {
		updateActiveNav();
	});
})();
