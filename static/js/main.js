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

		function showSuccess() {
			contactForm.outerHTML =
				'<div id="contact-form" class="form-success">' +
				'<div class="check">&#10003;</div>' +
				"<h3>Message sent!</h3>" +
				"<p>Thank you for reaching out. I&rsquo;ll get back to you within 24 hours.</p>" +
				"</div>";
		}

		if (fsId) {
			window.formspree = window.formspree || function () { (formspree.q = formspree.q || []).push(arguments); };
			formspree("initForm", { formElement: "#contact-form", formId: fsId });
			contactForm.addEventListener("formspree:success", function () {
				showSuccess();
			});
		} else if (window.location.hostname === "localhost" || window.location.hostname === "127.0.0.1") {
			var btn = contactForm.querySelector("button[type=submit]");
			contactForm.addEventListener("submit", function (e) {
				e.preventDefault();
				if (btn) {
					btn.disabled = true;
					btn.textContent = "Sending...";
				}
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
						if (btn) {
							btn.disabled = false;
							btn.textContent = "Send Message";
						}
						alert("Something went wrong. Please try again.");
					});
			});
		} else {
			contactForm.addEventListener("submit", function (e) {
				e.preventDefault();
				alert("Contact form not configured. Please email me directly at kylereynoldsdev@gmail.com.");
			});
		}
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
