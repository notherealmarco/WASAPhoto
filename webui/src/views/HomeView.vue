<script>
export default {
	data: function () {
		return {
			// Whether the content is loading
			// to show the loading spinner
			loading: false,

			// Stream data from the server
			stream_data: [],

			// Whether the data has ended
			// to stop loading more data with the infinite scroll
			data_ended: false,

			// Parameters to load data dynamically when scrolling
			start_idx: 0,
			limit: 1,

			// Shows the retry button
			loadingError: false,
		}
	},
	methods: {
		// Reload the whole page content
		// fetching it again from the server
		async refresh() {
			// Limits the number of posts to load based on the window height
			// to avoid loading too many posts at once
			// 450px is (a bit more) of the height of a single post
			this.limit = Math.round(window.innerHeight / 450);

			// Reset the parameters and the data
			this.start_idx = 0;
			this.data_ended = false;
			this.stream_data = [];

			// Fetch the first batch of posts
			this.loadContent();
		},

		// Requests data from the server asynchronously
		async loadContent() {
			this.loading = true;

			let response = await this.$axios.get("/stream?start_index=" + this.start_idx + "&limit=" + this.limit);

			// Errors are handled by the interceptor, which shows a modal dialog to the user and returns a null response.
			if (response == null) {
				this.loading = false
				this.loadingError = true
				return
			}

			// If the response is empty or shorter than the limit
			// then there is no more data to load
			if (response.data.length == 0 || response.data.length < this.limit) this.data_ended = true;
			this.stream_data = this.stream_data.concat(response.data);

			// Finished loading, hide the spinner
			this.loading = false;
		},

		// Loads more data when the user scrolls down
		// (this is called by the IntersectionObserver component)
		loadMore() {
			// Avoid loading more content if the data has ended
			if (this.loading || this.data_ended) return

			// Increase the start index and load more content
			this.start_idx += this.limit
			this.loadContent()
		},
	},

	// Called when the view is mounted
	mounted() {
		// Start loading the content
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

					<!-- Show a message if there's no content to show -->
					<div v-if="(stream_data.length == 0)" class="alert alert-secondary text-center" role="alert">
						There's nothing here 😢
						<br />Why don't you start following somebody? 👻
					</div>

					<!-- The stream -->
					<div id="main-content" v-for="item of stream_data" v-bind:key="item.photo_id">
						<!-- PostCard for each photo -->
						<PostCard :user_id="item.user_id" :photo_id="item.photo_id" :name="item.name" :date="item.date"
							:comments="item.comments" :likes="item.likes" :liked="item.liked" />
					</div>

					<!-- Show a message if there's no more content to show -->
					<div v-if="(data_ended && !(stream_data.length == 0))" class="alert alert-secondary text-center" role="alert">
						This is the end of your stream. Hooray! 👻
					</div>

					<!-- The loading spinner -->
					<LoadingSpinner :loading="loading" /><br />

					<div class="d-flex align-items-center flex-column">
						<!-- Retry button -->
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>

						<!-- Load more button -->
						<button v-if="(!data_ended && !loading)" @click="loadMore" class="btn btn-secondary py-1 mb-5"
							style="border-radius: 15px">Load more</button>

						<!-- The IntersectionObserver for dynamic loading -->
						<IntersectionObserver sentinal-name="load-more-home" @on-intersection-element="loadMore" />
					</div>
				</div>
			</div>
		</div>
	</div>
</template>