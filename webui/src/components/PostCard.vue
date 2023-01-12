<script>
export default {
	props: ["user_id", "name", "date", "comments", "likes", "photo_id", "liked"],
	data: function () {
		return {
			// Whether the image is loaded (disables the spinner)
			imageReady: false,

			// Likes and comments
			post_liked: this.liked,
			post_like_cnt: this.likes,
			post_comments_cnt: this.comments,
			comments_data: [],
			comments_start_idx: 0,
			comments_shown: false,
			commentMsg: "",

			// Whether the comments have ended (no more comments to load)
			data_ended: false,
		}
	},
	methods: {
		// Visit the user's profile
		visitUser() {
			this.$router.push({ path: "/profile/" + this.user_id });
		},

		// Post a new comment
		postComment() {
			this.$axios.post("/users/" + this.user_id + "/photos/" + this.photo_id + "/comments", {
				"comment": this.commentMsg,
				"user_id": this.$currentSession(),
			}).then(response => {
				if (response == null) return // the interceptors returns null if something goes bad

				// Reset the comment input and update the counter
				this.commentMsg = "";
				this.post_comments_cnt++;

				// Fetch comments from the server
				this.comments_data = [];
				this.comments_start_idx = 0;
				this.getComments();
			})
		},

		// Show or hide the comments section
		showHideComments() {
			// If comments are already shown, hide them and reset the data
			if (this.comments_shown) {
				this.comments_shown = false;
				this.comments_data = [];
				this.comments_start_idx = 0;
				return;
			}
			this.getComments();
		},

		// Fetch comments from the server
		getComments() {
			this.data_ended = false

			this.$axios.get("/users/" + this.user_id + "/photos/" + this.photo_id +
				"/comments?limit=5&start_index=" + this.comments_start_idx).then(response => {

					// If there are no more comments, set the flag
					if (response.data.length == 0) this.data_ended = true;

					// Otherwise increment the start index
					else this.comments_start_idx += 5;

					// Append the comments to the array (they will be rendered)
					this.comments_data = this.comments_data.concat(response.data);
					this.comments_shown = true;
				})
		},

		// Like the photo
		like() {
			this.$axios.put("/users/" + this.user_id + "/photos/" + this.photo_id + "/likes/" + this.$currentSession()).then(response => {
				if (response == null) return
				this.post_liked = true;
				this.post_like_cnt++;
			})
		},

		// Unlike the photo
		unlike() {
			this.$axios.delete("/users/" + this.user_id + "/photos/" + this.photo_id + "/likes/" + this.$currentSession()).then(response => {
				if (response == null) return
				this.post_liked = false;
				this.post_like_cnt--;
			})
		},
	},

	created() {
		// Fetch the image from the server and display it
		this.$axios.get("/users/" + this.user_id + "/photos/" + this.photo_id, {
			responseType: 'arraybuffer'
		}).then(response => {
			// Create an image element and append it to the container
			const img = document.createElement('img');

			// Set image source and css class
			img.src = URL.createObjectURL(new Blob([response.data]));
			img.classList.add("card-img-top");

			// Append the image to the container and disable the spinner
			this.$refs.imageContainer.appendChild(img);
			this.imageReady = true;
		});
	},
}
</script>

<template>
	<div class="card mb-5">
		<div ref="imageContainer">
			<div v-if="!imageReady" class="mt-3 mb-3">
				<LoadingSpinner :loading="!imageReady" />
			</div>
		</div>

		<div class="container">
			<div class="row">
				<div class="col-10">
					<div class="card-body">
						<h5 @click="visitUser" class="card-title d-inline-block" style="cursor: pointer">{{ name }}</h5>
						<p class="card-text">{{ new Date(Date.parse(date)) }}</p>
					</div>
				</div>

				<div class="col-2">
					<div class="card-body d-flex justify-content-end" style="display: inline-flex">
						<!-- not quite sure flex is the right property, but it works -->
						<a @click="showHideComments">
							<h5><i class="card-title bi bi-chat-right pe-1"></i></h5>
						</a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_comments_cnt }}</h6>
						<a v-if="!post_liked" @click="like">
							<h5><i class="card-title bi bi-suit-heart ps-2 pe-1 like-icon"></i></h5>
						</a>
						<a v-if="post_liked" @click="unlike">
							<h5><i class="card-title bi bi-heart-fill ps-2 pe-1 like-icon like-red"></i></h5>
						</a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_like_cnt }}</h6>
						<h5></h5>
					</div>
				</div>
			</div>
			<div v-if="comments_shown">
				<div v-for="item of comments_data" class="row" v-bind:key="item.comment_id">
					<div class="col-7 card-body border-top">
						<b>{{ item.name }}:</b> {{ item.comment }}
					</div>
					<div class="col-5 card-body border-top text-end text-secondary">
						{{ new Date(Date.parse(item.date)).toDateString() }}
					</div>
				</div>
				<div v-if="!data_ended" class="col-12 card-body text-end pt-0 pb-1 px-0">
					<a @click="getComments" class="text-primary">Mostra altro...</a>
				</div>
				<div class="row">
					<div class="col-10 card-body border-top text-end">
						<input v-model="commentMsg" type="text" class="form-control" placeholder="Commenta...">
					</div>
					<div class="col-1 card-body border-top text-end ps-0 d-flex">
						<button style="width: 100%" type="button" class="btn btn-primary"
							@click="postComment">Go</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.like-icon:hover {
	color: #ff0000;
}

.like-red {
	color: #ff0000;
}
</style>