<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			stream_data: null,
			my_id: sessionStorage.getItem("token"),
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/stream");
				this.stream_data = response.data;
				this.errormsg = this.stream_data; // TODO: temporary
				this.loading = false;
			} catch (e) {
				if (e.response.status == 401) {
					this.$router.push({ path: "/login" });
				}
				this.errormsg = e.toString();
			}
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

					<div id="main-content" v-for="item of stream_data">
						<PostCard :user_id="item.user_id"
									:photo_id="item.photo_id"
									:name="item.name"
									:date="item.date"
									:comments="item.comments"
									:likes="item.likes"
									:liked="item.liked"
									:my_id="my_id" />
					</div>

					<LoadingSpinner :loading="loading" /><br />
				</div>
			</div>
		</div>

	</div>
</template>

<style>
</style>
