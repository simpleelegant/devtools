// wrap all global variables and/or functions
window.A = {
    parseQueryString: function() {
        var qs = location.search.length ? location.search.substr(1).split('&') : [],
            args = {};

        qs.forEach(function(q) {
            if (q) {
                var kv = q.split('=');
                if (kv.length === 2) {
                    args[decodeURIComponent(kv[0])] = decodeURIComponent(kv[1]);
                }
            }
        });

        return args;
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

    request: function (method, url, data, success, fail) {
        $.ajax({
            method: method, // 'GET', 'POST', 'PUT', etc.
            url: url,
            data: data,
            success: success,
            error: function (jqXHR) {
                var msg = (jqXHR.responseJSON && jqXHR.responseJSON.error) ? 
                    jqXHR.responseJSON.error :
                    (jqXHR.responseText || jqXHR.status+' '+jqXHR.statusText);
                if (fail) {
                    fail(msg);
                } else {
                    alert(msg);
                }
			}
        });
    },

    api: {
        listJobs: function(server, namespace, labelSelector, success) {
            A.request('GET', '/kubernetes_dashboard/list_jobs', 
                {server: server, namespace: namespace, labelSelector: labelSelector}, success);
        },
        describeJob: function(server, namespace, name, success) {
            A.request('GET', '/kubernetes_dashboard/describe_job', 
                {server: server, namespace: namespace, name: name}, success);
        },
        deleteJob: function(server, namespace, name, success, fail) {
            var query = A.objectToQueryString({server: server, namespace: namespace, name: name});
            A.request('DELETE', '/kubernetes_dashboard/delete_job?'+query, 
                null, success, fail);
        },
        listPods: function(server, namespace, labelSelector, success) {
            A.request('GET', '/kubernetes_dashboard/list_pods', 
                {server: server, namespace: namespace, labelSelector: labelSelector}, success);
        },
        describePod: function(server, namespace, name, success) {
            A.request('GET', '/kubernetes_dashboard/describe_pod', 
                {server: server, namespace: namespace, name: name}, success);
        },
        deletePod: function(server, namespace, name, success) {
            var query = A.objectToQueryString({server: server, namespace: namespace, name: name});
            A.request('DELETE', '/kubernetes_dashboard/delete_pod?'+query, 
                null, success);
        }
    },

    listJobs: function(args) {
        var $panel = $('#list-jobs-panel'),
            $toolbar = $panel.find('.toolbar'),
            $cancel = $toolbar.find('.cancel'),
            $result = $toolbar.find('.result'),
            $body = $panel.find('tbody'),
            $row = $body.find('tr').detach(),
            interval = null,
            wantStopDeleting = false;

        $panel.show();

        $cancel.click(function() {
            wantStopDeleting = true;
            $cancel.hide();
            $result.text('Deleting cancel.').show();
        });

        $toolbar.find('.delete-jobs-btn').click(function() {
            if (!confirm("Sure?")) { return; }

            var $rows = $body.find('tr'),
                i = $rows.length - 1;

            $(this).hide();
            $cancel.show();

            interval = setInterval(function() {
                if (i < 0 || wantStopDeleting) {
                    clearInterval(interval);
                    interval = null;

                    if (i < 0) {
                        $cancel.hide();
                        $result.text('Deleting finished.').show();
                    }

                    return;
                }

                var $r = $rows.eq(i);

                i -= 1;

                $cancel.text($r.find('.order').text()+' job is deleting... Click me to cancel');

                A.api.deleteJob(args.server, $r.data('namespace'), $r.data('name'), function() {
                    $r.fadeOut(300, function() { $r.remove(); });
                }, function(err) {
                    // append error info if fails
                    $r.find('.delete').closest('td').append(err);
                });
            }, 1000); // wait 1 second for possible to cancel
        });

        A.api.listJobs(args.server, args.namespace, args.labelSelector, function(jobs) {
            $panel.find('.head-tip code').text((jobs && jobs.length) ? jobs.length : 0);
            if (!jobs || !jobs.length) {
                $body.append('<tr><td colspan="6" class="empty">No job.</td></tr>');
                return;
            }

            jobs.sort(function(i, j) {
                return i.creationTimestamp < j.creationTimestamp ? 1 : -1;
            });

            var param = { action: 'describeJob', server: args.server };

            jobs.forEach(function(job, i) {
                var $r = $row.clone();
                param.namespace = job.namespace;
                param.name = job.name;
                $r.attr('data-name', job.name);
                $r.attr('data-namespace', job.namespace);
                $r.find('.order').text('#'+(i+1));
                $r.find('.namespace').text(job.namespace);
                $r.find('.name').text(job.name);
                $r.find('.creationTimestamp').text(job.creationTimestamp);
                $r.find('.succeeded').text(job.succeeded || 0);
                $r.find('.detail').attr('href', '/kubernetes_dashboard/?'+A.objectToQueryString(param));
                $r.appendTo($body);
            });

            $toolbar.show();
        });
    },

    listPods: function(args) {
        var $panel = $('#list-pods-panel'),
            $toolbar = $panel.find('.toolbar'),
            $cancel = $toolbar.find('.cancel'),
            $result = $toolbar.find('.result'),
            $body = $panel.find('tbody'),
            $row = $body.find('tr').detach(),
            interval = null,
            wantStopDeleting = false;

        $panel.show();

        $cancel.click(function() {
            wantStopDeleting = true;
            $cancel.hide();
            $result.text('Deleting cancel.').show();
        });

        $toolbar.find('.delete-pods-btn').click(function() {
            if (!confirm("Sure?")) { return; }

            var $rows = $body.find('tr'),
                i = $rows.length - 1;

            $(this).hide();
            $cancel.show();

            interval = setInterval(function() {
                if (i < 0 || wantStopDeleting) {
                    clearInterval(interval);
                    interval = null;

                    if (i < 0) {
                        $cancel.hide();
                        $result.text('Deleting finished.').show();
                    }

                    return;
                }

                var $r = $rows.eq(i);

                i -= 1;

                $cancel.text($r.find('.order').text()+' pod is deleting... Click me to cancel');

                A.api.deletePod(args.server, $r.data('namespace'), $r.data('name'), function() {
                    $r.fadeOut(300, function() { $r.remove(); });
                }, function(err) {
                    // append error info if fails
                    $r.find('.delete').closest('td').append(err);
                });
            }, 1000); // wait 1 second for possible to cancel
        });

        A.api.listPods(args.server, args.namespace, args.labelSelector, function(pods) {
            $panel.find('.head-tip code').text((pods && pods.length) ? pods.length : 0);
            if (!pods || !pods.length) {
                $body.append('<tr><td colspan="6" class="empty">No pod.</td></tr>');
                return;
            }

            pods.sort(function(i, j) {
                return i.creationTimestamp < j.creationTimestamp ? 1 : -1;
            });

            var param = { action: 'describePod', server: args.server };

            pods.forEach(function(pod, i) {
                var $r = $row.clone();
                param.namespace = pod.namespace;
                param.name = pod.name;
                $r.attr('data-name', pod.name);
                $r.attr('data-namespace', pod.namespace);
                $r.find('.order').text('#'+(i+1));
                $r.find('.namespace').text(pod.namespace);
                $r.find('.name').text(pod.name);
                $r.find('.creationTimestamp').text(pod.creationTimestamp);
                $r.find('.node').text(pod.nodeName);
                $r.find('.detail').attr('href', '/kubernetes_dashboard/?'+A.objectToQueryString(param));
                $r.appendTo($body);
            });

            $toolbar.show();
        });
    },

    // show panel of Describe Job
    describeJob: function(args) {
        var $panel = $('#describe-job-panel'),
            $body = $panel.find('tbody'),
            $row = $body.find('tr').detach();

        $panel.show();

        $panel.find('.delete-job-btn').click(function() {
            if (!confirm("Sure?")) { return; }

            A.api.deleteJob(args.server, args.namespace, args.name, function() {
                $panel.html('<div class="empty">Job deleted.</div>');
            });
        });

        A.api.describeJob(args.server, args.namespace, args.name, function(data) {
            $panel.find('.toolbar').show();
            $panel.find('.head-tip code').text(args.name);
            $panel.find('pre code').text(data.job);

            // process pods
            if (!data.pods || !data.pods.items || data.pods.items.length===0) {
                $body.append('<tr><td colspan="5" class="empty">No pod.</td></tr>');
                return;
            }

            data.pods.items.sort(function(i, j) {
                return i.metadata.creationTimestamp < j.metadata.creationTimestamp ? 1 : -1;
            });

            var param = { action: 'describePod', server: args.server };

            data.pods.items.forEach(function(pod, i) {
                var $r = $row.clone();
                param.namespace = pod.metadata.namespace;
                param.name = pod.metadata.name;
                $r.attr('data-name', pod.metadata.name);
                $r.attr('data-namespace', pod.metadata.namespace);
                $r.find('.order').text('#'+(i+1));
                $r.find('.name').text(pod.metadata.name);
                $r.find('.creationTimestamp').text(pod.metadata.creationTimestamp);
                $r.find('.node').text(pod.spec.nodeName);
                $r.find('.detail').attr('href', '/kubernetes_dashboard/?'+A.objectToQueryString(param));
                $r.appendTo($body);
            });
        });
    },

    describePod: function(args) {
        var $panel = $('#describe-pod-panel');

        $panel.show();

        $panel.find('.delete-pod-btn').click(function() {
            if (!confirm("Sure?")) { return; }

            A.api.deletePod(args.server, args.namespace, args.name, function() {
                $panel.html('<div class="empty">Pod deleted.</div>');
            });
        });

        A.api.describePod(args.server, args.namespace, args.name, function(data) {
            $panel.find('.toolbar').show();
            $panel.find('.head-tip code').text(args.name);
            $('#pod-description').text(data.pod);
            $('#pod-logs pre code').text(data.logs);
        });
    },

    load: function() {
        var args = this.parseQueryString(),
            $server = $('input[name="server"]'),
            formatServer = function() {
                var s = $server.val().trim();

                if (!s.length) {
                    alert('kubernetes API Server must be not empty');
                    return false;
                }

                if (s.toLowerCase().indexOf('://') === -1) { $server.val('http://'+s); }

                return true;
            };

        $server.val(args.server || '');
        $('input[name="namespace"]').val(args.namespace || '');
        $('input[name="name"]').val(args.name || '');
        $('input[name="labelSelector"]').val(args.labelSelector || '');

        $server.focus();

        $('#search-jobs').click(function() {
            if (formatServer()) {
                location.href = '/kubernetes_dashboard/?action=listJobs&'+$('form').serialize();
            }
        });

        $('#search-pods').click(function() {
            if (formatServer()) {
                location.href = '/kubernetes_dashboard/?action=listPods&'+$('form').serialize();
            }
        });

        if (args.action) {
            this[args.action](args);
        }
    }
};

$(function() { A.load(); });
