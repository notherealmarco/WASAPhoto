<script>
export default {
	props: ["user_id", "name", "date", "comments", "likes", "photo_id", "liked"],
	data: function () {
		return {
			modalTitle: "Modal Title",
			modalMsg: "Modal Message",
		}
	},
	methods: {
		showModal(title, message) {
			this.modalTitle = title;
			this.modalMsg = message;

			// Simulate a click on the hidden modal button to open it
			this.$refs.openModal.click();
		},
	},

	mounted() {
		// Check if the user is already logged in
		this.$axiosUpdate()

		// Configure axios interceptors
		this.$axios.interceptors.response.use(response => {
			// Leave response as is
			return response;
		}, error => {
			if (error.response.status != 0) {
				// If the response is 401, redirect to /login
				if (error.response.status === 401) {
					this.$router.push({ path: '/login' })
					return;
				}
				
				// Show the error message from the server in a modal
				this.showModal("Error " + error.response.status, error.response.data['status'])
				return;
			}
			// Show the error message from axios in a modal
			this.showModal("Error", error.toString());
			return;
		});
	}
}
</script>

<template>
	<!-- Invisible button to open the modal -->
	<button ref="openModal" type="button" class="btn btn-primary" style="display: none" data-bs-toggle="modal" data-bs-target="#modal" />
	<!-- Modal to show error messages -->
	<Modal :title="modalTitle" :message="modalMsg" />

	<div class="container-fluid">
		<div class="row">
			<main class="mb-5">
				<!-- The view is rendered here -->
				<RouterView />
			</main>

			<!-- Bottom navigation buttons -->
			<nav id="global-nav" class="navbar fixed-bottom navbar-light bg-light row">
				<div class="collapse navbar-collapse" id="navbarNav"></div>
				<RouterLink to="/" class="col-4 text-center">
					<i class="bi bi-house text-dark" style="font-size: 2em"></i>
				</RouterLink>
				<RouterLink to="/search" class="col-4 text-center">
					<i class="bi bi-search text-dark" style="font-size: 2em"></i>
				</RouterLink>
				<RouterLink to="/profile/me" class="col-4 text-center">
					<i class="bi bi-person text-dark" style="font-size: 2em"></i>
				</RouterLink>
			</nav>
		</div>
	</div>
</template>

<style>
/* Make the active navigation button a little bit bigger */
#global-nav a.router-link-active {
	font-size: 1.2em
}
</style>
