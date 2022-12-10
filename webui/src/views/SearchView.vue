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
			field_username: "",
			my_id: sessionStorage.getItem("token"),
		}
	},
	methods: {
		async refresh() {
			this.limit = Math.round(window.innerHeight / 72);
			this.start_idx = 0;
			this.data_ended = false;
			this.stream_data = [];
			this.loadContent();
		},
		async loadContent() {
			this.loading = true;
			this.errormsg = null;
			if (this.field_username == "") {
				this.errormsg = "Please enter a username";
				this.loading = false;
				return;
			}
			try {
				let response = await this.$axios.get("/users?query=" + this.field_username + "&start_index=" + this.start_idx + "&limit=" + this.limit);
				if (response.data.length == 0) this.data_ended = true;
				else this.stream_data = this.stream_data.concat(response.data);
				this.loading = false;
			} catch (e) {
				this.errormsg = e.toString();
				if (e.response.status == 401) {
					this.$router.push({ path: "/login" });
				}
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
		// 72 is a bit more of the max height of a card
		// todo: may not work in 4k screens :/
		this.limit = Math.round(window.innerHeight / 72);
		this.scroll();
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
                        <input v-model="field_username" @input="refresh" id="formUsername" class="form-control" placeholder="name@example.com"/>
                        <label class="form-label" for="formUsername">Search by username</label>
                    </div>

					<div id="main-content" v-for="item of stream_data">
						<UserCard
								:user_id="item.user_id"
								:name="item.name"
								:followed="item.followed"
								:banned="item.banned"
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
