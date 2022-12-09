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

            user_data: [],
		}
	},
	methods: {
		async refresh() {
            this.getMainData();
			this.limit = Math.round(window.innerHeight / 450);
			this.start_idx = 0;
			this.data_ended = false;
			this.stream_data = [];
			this.loadContent();
		},


        async getMainData() {
            try {
                let response = await this.$axios.get("/users/" + this.$route.params.user_id);
                this.user_data = response.data;
            } catch(e) {
                this.errormsg = e.toString();
            }
        },

		async loadContent() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + this.$route.params.user_id + "/photos" + "?start_index=" + this.start_idx + "&limit=" + this.limit);
				if (response.data.length == 0) this.data_ended = true;
				else this.stream_data = this.stream_data.concat(response.data);
				this.loading = false;
			} catch (e) {
				if (e.response.status == 401) { // todo: move from here
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
        this.getMainData();
		this.limit = Math.round(window.innerHeight / 450);
		this.scroll();
		this.loadContent();
	}
}
</script>

<template>
	<div class="mt-5">

		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

                    <UserCard :user_id = "$route.params.user_id"
                                :name = "user_data['name']"
                                :followed = "user_data['followed']"
                                :banned = "user_data['banned']"
                                :my_id = "my_id"
                                :show_new_post = "true"
                                @updateInfo = "getMainData" />
                    
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
						<PostCard :user_id = "$route.params.user_id"
									:photo_id = "item.photo_id"
									:name = "user_data['name']"
									:date = "item.date"
									:comments = "item.comments"
									:likes = "item.likes"
									:liked = "item.liked"
									:my_id = "my_id" />
					</div>

					<div v-if="data_ended" class="alert alert-secondary text-center" role="alert">
						Hai visualizzato tutti i post. Hooray! 👻
					</div>

					<LoadingSpinner :loading="loading" /><br />
				</div>
			</div>
		</div>

	</div>
</template>

<style>
</style>
