<script>
// import getCurrentSession from '../services/authentication'; todo: can be removed
export default {
	data: function () {
		return {
			// The error message to display
			errormsg: null,

			loading: false,

			// Search results
			streamData: [],

			// Dynamic loading
			dataEnded: false,
			startIdx: 0,
			limit: 1,

			// Search input
			fieldUsername: "",
		}
	},
	methods: {
		// Reset the results and fetch the new requested ones
		async query() {
			// Set the limit to the number of cards that can fit in the window
			this.limit = Math.round(window.innerHeight / 72);

			// Reset the parameters and the data
			this.startIdx = 0;
			this.dataEnded = false;
			this.streamData = [];

			// Fetch the first batch of results
			this.loadContent();
		},

		// Fetch the search results from the server
		async loadContent() {
			this.loading = true;
			this.errormsg = null;

			// Check if the username is empty
			// and show an error message
			if (this.fieldUsername == "") {
				this.errormsg = "Please enter a username";
				this.loading = false;
				return;
			}

			// Fetch the results from the server
			let response = await this.$axios.get("/users?query=" + this.fieldUsername + "&start_index=" + this.startIdx + "&limit=" + this.limit);

			// Errors are handled by the interceptor, which shows a modal dialog to the user and returns a null response.
			if (response == null) {
				this.loading = false
				return
			}

			// If there are no more results, set the dataEnded flag
			if (response.data.length == 0) this.dataEnded = true;

			// Otherwise, append the new results to the array
			else this.streamData = this.streamData.concat(response.data);

			// Hide the loading spinner
			this.loading = false;
		},

		// Load a new batch of results when the user scrolls to the bottom of the page
		loadMore() {
			if (this.loading || this.dataEnded) return
			this.startIdx += this.limit
			this.loadContent()
		},
	},
}
</script>

<template>
	<div class="mt-4">
		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<h3 class="card-title border-bottom mb-4 pb-2 text-center">WASASearch</h3>

					<!-- Error message -->
					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

					<!-- Search form -->
					<div class="form-floating mb-4">
						<input v-model="fieldUsername" @input="query" id="formUsername" class="form-control"
							placeholder="name@example.com" />
						<label class="form-label" for="formUsername">Search by username</label>
					</div>

					<!-- Search results -->
					<div id="main-content" v-for="item of streamData" v-bind:key="item.user_id">
						<!-- User card (search result entry) -->
						<UserCard :user_id="item.user_id" :name="item.name" :followed="item.followed"
							:banned="item.banned" />
					</div>

					<!-- Loading spinner -->
					<LoadingSpinner :loading="loading" /><br />

					<!-- The IntersectionObserver for dynamic loading -->
					<IntersectionObserver sentinal-name="load-more-search" @on-intersection-element="loadMore" />
				</div>
			</div>
		</div>

	</div>
</template>

<style>

</style>
