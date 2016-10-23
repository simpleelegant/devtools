$(function() {
    window.custom = {
        serverFromQueryString: '',

        getParameterByName: function(name, url) {
            if (!url) url = window.location.href;
            name = name.replace(/[\[\]]/g, "\\$&");
            var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
                results = regex.exec(url);

            if (!results) return ''; // the name not exists
            if (!results[2]) return '';

            return decodeURIComponent(results[2].replace(/\+/g, " "));
        },

        objectToQueryString: function(obj) {
            var s = '';
            for (var key in obj) {
                if (obj.hasOwnProperty(key)) {
                    s += '&' + encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]);
                }
            }

            if (s.length && s[0] === '&') {
                s = s.substring(1);
            }

            return s;
        },

        formatDockerTimestamp: function(timestamp) {
            var t = new Date(timestamp * 1000);

            return t.toISOString();
        },

        toHumanReadableSize: function(bytes) {
            if (bytes >= 1024 * 1024) {
                return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
            }
            if (bytes >= 1024) {
                return (bytes / 1024).toFixed(1) + 'KB';
            }
            return bytes + 'B';
        },

        // inspect an image
        inspectImage: function(btn) {
            var $row = $(btn).closest('tr'),
                param = {
                    action: 'inspect_image',
                    id: $row.data('id'),
                    docker_api_server: this.serverFromQueryString
                };

            window.open('/docker_dashboard/?'+this.objectToQueryString(param));
        },

        // delete image
        deleteImage: function(btn) {
            var $row = $(btn).closest('tr');
            
            if (!confirm("Sure?")) { return; }

            $.post('/docker_dashboard/delete_image', {
                id: $row.data('id'),
                docker_api_server: this.serverFromQueryString
            }, function(data) {
                if (data.Error) { alert(data.Error); return; }
                $row.fadeOut(300, function() { $row.remove(); });
            });
        },

        // show panel of image getting result
        showPanelOfGetImagesResult: function() {
            var $panel = $('#images-panel'),
                $body = $panel.find('tbody'),
                $row = $body.find('tr').detach(),
                param = {
                    docker_api_server: this.serverFromQueryString
                };

            $panel.show();

            $.post('/docker_dashboard/get_images', param, function(data) {
                if (data.Error) { alert(data.Error); return; }
                if (!data.Data || !Array.isArray(data.Data)) { alert('No data.'); return; }

                var images = data.Data;

                $panel.find('.head-tip code').text(images.length);
                if (images.length===0) {
                    $body.append('<tr><td colspan="6" class="empty">No image.</td></tr>');
                    return;
                }

                images.sort(function(i, j) {
                    return i.Created < j.Created ? 1 : -1;
                });
                images.forEach(function(image, i) {
                    var $r = $row.clone();
                    $r.attr('data-id', image.Id);
                    $r.find('.order').text('#'+(i+1));
                    $r.find('.image-id').text(image.Id);
                    $r.find('.created').text(custom.formatDockerTimestamp(image.Created));
                    $r.find('.size').text(custom.toHumanReadableSize(image.Size));
                    image.RepoTags.forEach(function(t) {
                        $r.find('.repo-tags').append($('<div></div>').text(t));
                    });
                    $r.appendTo($body);
                });
            });
        },

        showPanelOfInspectImage: function() {
            var $panel = $('#image-panel'),
                param = {
                    id: this.getParameterByName('id'),
                    docker_api_server: this.serverFromQueryString
                };

            $panel.show();

            $panel.find('.delete').click(function() {
                if (!confirm("Sure?")) { return; }

                $.post('/docker_dashboard/delete_image', param, function(data) {
                    if (data.Error) { alert(data.Error); return; }
                    $panel.html('<div class="empty">Image deleted.</div>');
                });
            });

            $.post('/docker_dashboard/inspect_image', param, function(data) {
                if (data.Error) { alert(data.Error); return; }
                if (!data.Data) { alert('No data.'); return; }

                $panel.find('.toolbar').show();
                $panel.find('.head-tip code').text(param.id);
                $panel.find('pre code').text(data.Data);
            });
        },

        loadPage: function() {
            var $server = $('input[name="docker_api_server"]');

            this.serverFromQueryString = this.getParameterByName('docker_api_server');
            $server.val(this.serverFromQueryString);

            $server.focus();

            $('form').submit(function() { return false; }); 

            $('#get-images').click(function() {
                var server = $server.val().trim();
                if (server.length === 0) { alert('Docker API Server must be not empty'); return; }
                if (server.toLowerCase().indexOf('://') === -1) { $server.val('http://'+server); }
                location.href = '/docker_dashboard/?action=get_images_result&'+$('form').serialize();
            });

            switch (this.getParameterByName('action')) {
                case '':
                    break;
                case 'get_images_result':
                    this.showPanelOfGetImagesResult();
                    break;
                case 'inspect_image':
                    this.showPanelOfInspectImage();
                    break;
                default:
                    alert('Unkown action');
            }
        }
    };

    custom.loadPage();
});
