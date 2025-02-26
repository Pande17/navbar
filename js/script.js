document.getElementById('loginBtn').addEventListener('click', function () {
	document.getElementById('loginForm').style.display = 'block';
});

document.getElementById('registerBtn').addEventListener('click', function (event) {
	event.preventDefault();
	document.getElementById('registerForm').style.display = 'block';
	document.getElementById('loginForm').style.display = 'none';
});

document.querySelector('a[href="#loginForm"]').addEventListener('click', function (event) {
	event.preventDefault();
	document.getElementById('loginForm').style.display = 'block';
	document.getElementById('registerForm').style.display = 'none';
});

document.querySelectorAll('.close-btn').forEach(function (btn) {
	btn.addEventListener('click', function () {
		document.getElementById('loginForm').style.display = 'none';
		document.getElementById('registerForm').style.display = 'none';
	});
});

window.onclick = function (event) {
	if (event.target == document.getElementById('loginForm')) {
		document.getElementById('loginForm').style.display = 'none';
	}
	if (event.target == document.getElementById('registerForm')) {
		document.getElementById('registerForm').style.display = 'none';
	}
};
