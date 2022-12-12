<script setup>
import { RouterLink, RouterView } from 'vue-router'
import getCurrentSession from './services/authentication';
import { updateToken } from './services/axios';
</script>

<template>

	<!--<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAPhoto</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>-->

	<div class="container-fluid">
		<div class="row">
			<!--<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>WASAPhoto</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Stream
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link1" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Search
							</RouterLink>
						</li>
					</ul>

					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Account</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
								Your profile
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>-->

			<!---<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">-->
			<main class="mb-5">
				<RouterView />
			</main>

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
#global-nav a.router-link-active {
	font-size: 1.2em
}
</style>
