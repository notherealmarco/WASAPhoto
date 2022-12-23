<script>
import IntersectionObserver from '../components/IntersectionObserver.vue';

export default {
	data: function () {
		return {
			requestedProfile: this.$route.params.user_id,
			loading: true,
			loadingError: false,
			stream_data: [],
			data_ended: false,
			start_idx: 0,
			limit: 1,
			user_data: [],
		};
	},
	methods: {
		async refresh() {
			this.getMainData();
			// this way we are sure that we fill the first page todo: check
			// 450 is a bit more of the max height of a post
			// todo: may not work in 4k screens :/
			this.limit = Math.max(Math.round(window.innerHeight / 450), 1);
			this.start_idx = 0;
			this.data_ended = false;
			this.stream_data = [];
			this.loadContent();
		},
		async getMainData() {
			let response = await this.$axios.get("/users/" + this.requestedProfile);
			if (response == null) {
				this.loading = false;
				this.loadingError = true;
				return;
			}
			this.user_data = response.data;
		},
		async loadContent() {
			this.loading = true;
			let response = await this.$axios.get("/users/" + this.requestedProfile + "/photos" + "?start_index=" + this.start_idx + "&limit=" + this.limit);
			if (response == null) {
				// do something
				return;
			}
			if (response.data.length == 0 || response.data.length < this.limit)
				this.data_ended = true;
			this.stream_data = this.stream_data.concat(response.data);
			this.loading = false;
		},
		loadMore() {
			if (this.loading || this.data_ended) return
			this.start_idx += this.limit
			this.loadContent()
		},
	},
	created() {
		if (this.$route.params.user_id == "me") {
			//this.$router.replace({ path: "/profile/" +  }); (It's ok to not redirect, it's just a matter of taste)
			this.requestedProfile = this.$currentSession();
		}
		else {
			this.requestedProfile = this.$route.params.user_id;
		}
		//this.scroll();
		this.refresh();
	},
	components: { IntersectionObserver }
}
</script>

<template>
	<div class="mt-5">

		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<UserCard :user_id="requestedProfile" :name="user_data['name']" :followed="user_data['followed']"
						:banned="user_data['banned']" :my_id="this.$currentSession" :show_new_post="true"
						@updateInfo="getMainData" @updatePosts="refresh" />

					<div class="row text-center mt-2 mb-3">
						<div class="col-4" style="border-right: 1px">
							<h3>{{ user_data["photos"] }}</h3>
							<h6>Photos</h6>
						</div>
						<div class="col-4">
							<h3>{{ user_data["followers"] }}</h3>
							<h6>Followers</h6>
						</div>
						<div class="col-4">
							<h3>{{ user_data["following"] }}</h3>
							<h6>Following</h6>
						</div>
					</div>

					<div id="main-content" v-for="item of stream_data">
						<PostCard :user_id="requestedProfile" :photo_id="item.photo_id" :name="user_data['name']"
							:date="item.date" :comments="item.comments" :likes="item.likes" :liked="item.liked" />
					</div>

					<div v-if="data_ended" class="alert alert-secondary text-center" role="alert">
						You reached the end. Hooray! 👻
					</div>

					<LoadingSpinner :loading="loading" />

					<div class="d-flex align-items-center flex-column">
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>

						<button v-if="(!data_ended && !loading)" @click="loadMore" class="btn btn-secondary py-1 mb-5"
							style="border-radius: 15px">Load more</button>
						<IntersectionObserver sentinal-name="load-more-profile" @on-intersection-element="loadMore" />
					</div>
				</div>
			</div>
		</div>
	</div>
</template>