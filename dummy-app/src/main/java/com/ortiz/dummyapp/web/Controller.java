package com.ortiz.dummyapp.web;

import com.amazonaws.services.simplesystemsmanagement.AWSSimpleSystemsManagement;
import com.ortiz.dummyapp.domains.Post;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

@RestController
public class Controller {

  @Autowired
  private AWSSimpleSystemsManagement ssmClient;

  private HashMap<Long, Post> posts = new HashMap<>();

  private AtomicLong postGenId = new AtomicLong();

  @GetMapping("/posts")
  public List<Post> getParameter() {
    return new ArrayList<>(posts.values());
  }

  @GetMapping("/posts/{id}")
  public Post getPostById(@PathVariable Long id) {
    Post post = posts.get(id);
    if (post == null) {
      throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Post not found.");
    }
    return post;
  }

  @PostMapping("/posts")
  public Post createPost(Post post) {
    Long id = postGenId.incrementAndGet();
    post.setId(id);
    posts.put(id, post);
    return post;
  }

  @PutMapping("/posts/{id}")
  public Post updatePost(@PathVariable Long id, Post post) {
    Post postDb = posts.get(id);
    if (postDb == null) {
      throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Post not found.");
    }
    post.setId(id);
    posts.put(id, post);
    return post;
  }

  @PatchMapping("/posts/{id}")
  public Post updatePartially(@PathVariable Long id, Post post) {
    return updatePost(id, post);
  }

  @DeleteMapping("/posts/{id}")
  public void deletePost(@PathVariable Long id) {
    Post post = posts.get(id);
    if (post == null) {
      throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Post not found.");
    }
    posts.remove(id);
  }


}
