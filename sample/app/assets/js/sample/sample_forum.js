function initForumView() {
    console.log("initForumView");
	$("#forum-filter-button").on("click", handle_forum_filter_button_click);
    $("#forum-reply-button").on("click", handle_forum_reply_button_click);
    $("#forum-new-topic-button").on("click", handle_forum_newtopic_button_click);

}

function handle_forum_filter_button_click(e) {
	// dev debug printing
    console.log("filter button was clicked");
    console.log(e);

    // post filter values to server
    var id = $("#forum-filter-topicid").val();
    var csrf = $("#csrf_token").val();

    console.log("id", id);

    // actually send data
    dosend_forum_filter_update(id, csrf);
}

function handle_forum_reply_button_click(e) {
    // dev debug printing
    console.log("reply button was clicked");
    console.log(e);

       var forum_newmessage_text = [
    '    <div class="row">',
    '       <div class="large-12 small-12 columns">',
    '           <input type="hidden" name="csrf_token" value="{{ .csrf_token }}" />',
    '           <a class="button" id="message_post_button">Post</a>',
    '       </div>',
    '    </div>',
    '    <div class="row">',
    '       <div id="epiceditor"></div>',
    '    </div>',
    ].join("\n");


    $("#forum-topic-new-message").append(forum_newmessage_text);

    $("#message_post_button").on("click", handle_newmessage_post_button_click);

    initReplyEditor();

}

function handle_forum_newtopic_button_click(e) {
    // dev debug printing
    console.log("new  topic button was clicked");
    console.log(e);
    var forum_newtopic_text = [
    '    <div class="row">',
    '       <div class="large-9 small-9 columns">',
    '           <input type="hidden" name="csrf_token" value="{{ .csrf_token }}" />',
    '           <input id="subject_field" type="text" placeholder="Subject">',
    '       </div>',
    '       <div class="large-3 small-3 columns">',
    '           <a class="button" id="topic_post_button">Post</a>',
    '       </div>',
    '    </div>',
    '    <div class="row">',
    '       <div id="epiceditor"></div>',
    '    </div>',
    ].join("\n");


    $("#forum-topiclist-new-topic").append(forum_newtopic_text);

    $("#topic_post_button").on("click", handle_newtopic_post_button_click);

    initTopicEditor();

}

function handle_newmessage_post_button_click(e) {
    console.log("new message post button was clicked");

    var topicId = $("#forum-topicId").text();

    var content = editor.exportFile();
    console.log(content);
    console.log(topicId);

    var csrf = $("#csrf_token").val();
    dosend_newmessage_post(content, topicId, csrf);
}

function handle_newtopic_post_button_click(e) {
    console.log("new topic post button was clicked");

    var subject = $("#subject_field").val();

    var content = editor.exportFile();
    console.log(content);

    var csrf = $("#csrf_token").val();
    dosend_newtopic_post(subject, content, csrf);
}

function dosend_forum_filter_update(id, csrf) {

    var post_query = "/forum/filter";
    post_query += "?id=" + encodeURIComponent(id);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", post_query, true);
    xhr.setRequestHeader('X-CSRF-Token', csrf);


    xhr.onload = function(e) {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var results = JSON.parse(xhr.responseText);
                console.log(results);
                update_forum_topiclist_results_table(results);
            } else {
                console.error(xhr.statusText);
            }
        }
    };



    xhr.onerror = function(e) {
        console.error(xhr.statusText);
    };
    xhr.send(null);
}

function dosend_newmessage_post(content, topicId, csrf) {
    var post_query = "/forum/message";

    post_query += "?topicId=" + encodeURIComponent(topicId);
    post_query += "&content=" + encodeURIComponent(content);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", post_query, true);
    xhr.setRequestHeader('X-CSRF-Token', csrf);


    xhr.onload = function(e) {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                // var results = JSON.parse(xhr.responseText);
                console.log(xhr.responseText);
                
            } else {
                console.error(xhr.statusText);
            }
        }
    };



    xhr.onerror = function(e) {
        console.error(xhr.statusText);
    };
    xhr.send(null);
}

function dosend_newtopic_post(subject, content, csrf) {
    var post_query = "/forum/topic";
    post_query += "?subject=" + encodeURIComponent(subject);
    post_query += "&content=" + encodeURIComponent(content);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", post_query, true);
    xhr.setRequestHeader('X-CSRF-Token', csrf);


    xhr.onload = function(e) {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var results = JSON.parse(xhr.responseText);
                console.log(xhr.responseText);
                update_forum_topiclist_results_table(results);
            } else {
                console.error(xhr.statusText);
            }
        }
    };



    xhr.onerror = function(e) {
        console.error(xhr.statusText);
    };
    xhr.send(null);
}

function update_forum_topiclist_results_table(results) {
    // clear the current table results
    $("#forum-topiclist-results").empty();

    // render results in panel
    if (results === null) {
        return
    }

    var template = Hogan.compile(forum_topic_row_template_text, { delimiters: '<% %>' });
    if (results instanceof Array) {
        for (var i = 0; i < results.length; i++) {
            var output = template.render(results[i]);
            $("#forum-topiclist-results").append(output)
        }
    } else {  // should be a single element
            var output = template.render(results);
            $("#forum-topiclist-results").append(output)
    }


}

var forum_topic_row_template_text = [
'<tbody>',
'   <tr>',
'       <th>{{.TopicId}}</th>',
'       <th>',
'           <a href="/forum/{{.TopicId}}">{{.TopicName}}</a>',
'       </th>',
'       <th>{{.TopicTags}}</th>',
'       <th>{{.OpenedBy}}</th>',
'       <th>{{.OpenedAt}}</th>',
'   </tr>',
'</tbody>',
].join("\n");
