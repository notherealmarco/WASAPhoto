<script>
export default {
	data: function () {
		return {
			loading: false,
			stream_data: [],
			data_ended: false,
			start_idx: 0,
			limit: 1,
			loadingError: false,
		}
	},
	methods: {
		async refresh() {
			// this way we are sure that we fill the first page
			// 450 is a bit more of the max height of a post
			// todo: may not work in 4k screens :/
			this.limit = Math.round(window.innerHeight / 450);
			this.start_idx = 0;
			this.data_ended = false;
			this.stream_data = [];
			this.loadContent();
		},
		async loadContent() {
			this.loading = true;

			let response = await this.$axios.get("/stream?start_index=" + this.start_idx + "&limit=" + this.limit);

			// Errors are handled by the interceptor, which shows a modal dialog to the user and returns a null response.
			if (response == null) {
				this.loading = false
				this.loadingError = true
				return
			}

			if (response.data.length == 0 || response.data.length < this.limit) this.data_ended = true;
			this.stream_data = this.stream_data.concat(response.data);
			this.loading = false;
		},
		loadMore() {
			if (this.loading || this.data_ended) return
			this.start_idx += this.limit
			this.loadContent()
		},
	},
	mounted() {
		this.refresh();
	}
}
</script>

<template>
	<div class="mt-4">
		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">
					<h3 class="card-title border-bottom mb-4 pb-2 text-center">Your daily WASAStream!</h3>

					<div v-if="(stream_data.length == 0)" class="alert alert-secondary text-center" role="alert">
						There's nothing here 😢
						<br />Why don't you start following somebody? 👻
					</div>

					<div id="main-content" v-for="item of stream_data">
						<PostCard :user_id="item.user_id" :photo_id="item.photo_id" :name="item.name" :date="item.date"
							:comments="item.comments" :likes="item.likes" :liked="item.liked" />
					</div>

					<div v-if="data_ended" class="alert alert-secondary text-center" role="alert">
						This is the end of your stream. Hooray! 👻
					</div>

					<LoadingSpinner :loading="loading" /><br />

					<div class="d-flex align-items-center flex-column">
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>

						<button v-if="(!data_ended && !loading)" @click="loadMore" class="btn btn-secondary py-1 mb-5"
							style="border-radius: 15px">Load more</button>
						<IntersectionObserver sentinal-name="load-more-home" @on-intersection-element="loadMore" />
					</div>
				</div>
			</div>
		</div>
	</div>
</template>