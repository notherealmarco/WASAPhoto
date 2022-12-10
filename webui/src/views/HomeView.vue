<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			stream_data: [],
			data_ended: false,
			start_idx: 0,
			limit: 1,
			my_id: sessionStorage.getItem("token"),
		}
	},
	methods: {
		async refresh() {
			this.limit = Math.round(window.innerHeight / 450);
			this.start_idx = 0;
			this.data_ended = false;
			this.stream_data = [];
			this.loadContent();
		},
		async loadContent() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/stream?start_index=" + this.start_idx + "&limit=" + this.limit);
				if (response.data.length == 0) this.data_ended = true;
				else this.stream_data = this.stream_data.concat(response.data);
				this.loading = false;
			} catch (e) {
				if (e.response.status == 401) {
					this.$router.push({ path: "/login" });
				}
				this.errormsg = e.toString();
			}
		},
		scroll () {
			window.onscroll = () => {
				let bottomOfWindow = Math.max(window.pageYOffset, document.documentElement.scrollTop, document.body.scrollTop) + window.innerHeight === document.documentElement.offsetHeight
				if (bottomOfWindow && !this.data_ended) {
					this.start_idx += this.limit;
					this.loadContent();
				}
			}
		},
	},
	mounted() {
		// this way we are sure that we fill the first page
		// 450 is a bit more of the max height of a post
		// todo: may not work in 4k screens :/
		this.limit = Math.round(window.innerHeight / 450);
		this.scroll();
		this.loadContent();
	}
}
</script>

<template>
	<div class="mt-4">
		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">
					<h3 class="card-title border-bottom mb-4 pb-2 text-center">Your daily WASAStream!</h3>

					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

					<div v-if="(stream_data.length == 0)" class="alert alert-secondary text-center" role="alert">
						There's nothing here 😢
						<br />Why don't you start following somebody? 👻
					</div>

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

					<div v-if="data_ended" class="alert alert-secondary text-center" role="alert">
						This is the end of your stream. Hooray! 👻
					</div>

					<LoadingSpinner :loading="loading" /><br />
				</div>
			</div>
		</div>

	</div>
</template>

<style>
</style>
