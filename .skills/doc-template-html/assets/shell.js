(function () {
  var bar = document.getElementById("progress");
  var links = {};
  document.querySelectorAll(".side nav a").forEach(function (a) {
    var href = a.getAttribute("href");
    if (href && href.charAt(0) === "#") links[href.slice(1)] = a;
  });
  var visible = new Set();
  function onScroll() {
    var h = document.documentElement;
    var max = h.scrollHeight - h.clientHeight;
    if (bar && max > 0) bar.style.transform = "scaleX(" + (h.scrollTop / max) + ")";
  }
  window.addEventListener("scroll", onScroll, { passive: true });
  onScroll();
  if (!("IntersectionObserver" in window)) return;
  var io = new IntersectionObserver(function (entries) {
    entries.forEach(function (e) {
      visible[e.isIntersecting ? "add" : "delete"](e.target.id);
    });
    var first = Object.keys(links).find(function (id) { return visible.has(id); });
    Object.keys(links).forEach(function (id) {
      links[id].classList.toggle("active", id === first);
    });
  }, { rootMargin: "0px 0px -70% 0px" });
  document.querySelectorAll("h2[id]").forEach(function (h) { io.observe(h); });
})();
