<script>
// import getCurrentSession from '../services/authentication'; todo: can be removed
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			streamData: [],
			dataEnded: false,
			startIdx: 0,
			limit: 1,
			fieldUsername: "",
		}
	},
	methods: {
		async refresh() {
			this.limit = Math.round(window.innerHeight / 72);
			this.startIdx = 0;
			this.dataEnded = false;
			this.streamData = [];
			this.loadContent();
		},
		async loadContent() {
			this.loading = true;
			this.errormsg = null;
			if (this.fieldUsername == "") {
				this.errormsg = "Please enter a username";
				this.loading = false;
				return;
			}

			let response = await this.$axios.get("/users?query=" + this.fieldUsername + "&start_index=" + this.startIdx + "&limit=" + this.limit);

			// Errors are handled by the interceptor, which shows a modal dialog to the user and returns a null response.
			if (response == null) {
				this.loading = false
				return
			}

			if (response.data.length == 0) this.dataEnded = true;
			else this.streamData = this.streamData.concat(response.data);
			this.loading = false;

		},
		loadMore() {
			if (this.loading || this.dataEnded) return
			this.startIdx += this.limit
			this.loadContent()
		},
	},
	mounted() {
		// this way we are sure that we fill the first page
		// 72 is a bit more of the max height of a card
		// todo: may not work in 4k screens :/
		this.limit = Math.round(window.innerHeight / 72)
	}
}
</script>

<template>
	<div class="mt-4">
		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<h3 class="card-title border-bottom mb-4 pb-2 text-center">WASASearch</h3>

					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

					<div class="form-floating mb-4">
						<input v-model="fieldUsername" @input="refresh" id="formUsername" class="form-control"
							placeholder="name@example.com" />
						<label class="form-label" for="formUsername">Search by username</label>
					</div>

					<div id="main-content" v-for="item of streamData">
						<UserCard :user_id="item.user_id" :name="item.name" :followed="item.followed"
							:banned="item.banned" />
					</div>

					<LoadingSpinner :loading="loading" /><br />
					<IntersectionObserver sentinal-name="load-more-search" @on-intersection-element="loadMore" />
				</div>
			</div>
		</div>

	</div>
</template>

<style>

</style>
